package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend/internal/domain/course"
	"backend/internal/domain/grade"
	"backend/internal/domain/user"
)

type GradeService struct {
	gradeRepo  grade.Repository
	userRepo   user.Repository
	courseRepo course.Repository
}

func NewGradeService(
	gradeRepo grade.Repository,
	userRepo user.Repository,
	courseRepo course.Repository,
) *GradeService {
	return &GradeService{
		gradeRepo:  gradeRepo,
		userRepo:   userRepo,
		courseRepo: courseRepo,
	}
}

// GetMyGrades получает все оценки текущего студента
func (s *GradeService) GetMyGrades(ctx context.Context, userID uuid.UUID) ([]*grade.GradeResponse, error) {
	// Получаем оценки
	grades, err := s.gradeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Создаём карту курсов для быстрого доступа к названиям
	courseNames := make(map[uuid.UUID]string)

	// Преобразуем в Response
	responses := make([]*grade.GradeResponse, len(grades))
	for i, g := range grades {
		// Получаем название курса (если ещё не получали)
		if _, ok := courseNames[g.CourseID]; !ok {
			c, err := s.courseRepo.GetByID(ctx, g.CourseID)
			if err != nil || c == nil {
				courseNames[g.CourseID] = "Неизвестный курс"
			} else {
				courseNames[g.CourseID] = c.Name
			}
		}

		responses[i] = g.ToResponse(courseNames[g.CourseID])
	}

	return responses, nil
}

// GetMyGradesByCourse получает оценки по конкретному курсу
func (s *GradeService) GetMyGradesByCourse(ctx context.Context, userID, courseID uuid.UUID) ([]*grade.GradeResponse, error) {
	// Проверяем существование курса
	c, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("course not found")
	}

	// Получаем оценки
	grades, err := s.gradeRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	// Преобразуем в Response
	responses := make([]*grade.GradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = g.ToResponse(c.Name)
	}

	return responses, nil
}

// GetMyGradesByPeriod получает оценки за период
func (s *GradeService) GetMyGradesByPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*grade.GradeResponse, error) {
	// Валидация дат
	if from.After(to) {
		return nil, errors.New("from date must be before to date")
	}

	// Получаем оценки
	grades, err := s.gradeRepo.GetByUserAndPeriod(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	// Создаём карту курсов
	courseNames := make(map[uuid.UUID]string)

	// Преобразуем в Response
	responses := make([]*grade.GradeResponse, len(grades))
	for i, g := range grades {
		if _, ok := courseNames[g.CourseID]; !ok {
			c, err := s.courseRepo.GetByID(ctx, g.CourseID)
			if err != nil || c == nil {
				courseNames[g.CourseID] = "Неизвестный курс"
			} else {
				courseNames[g.CourseID] = c.Name
			}
		}

		responses[i] = g.ToResponse(courseNames[g.CourseID])
	}

	return responses, nil
}

// GetMySummary получает сводку успеваемости
func (s *GradeService) GetMySummary(ctx context.Context, userID uuid.UUID) (*grade.StudentSummary, error) {
	return s.gradeRepo.GetSummaryByUser(ctx, userID)
}

// CreateGrade создаёт новую оценку (для преподавателя/админа)
func (s *GradeService) CreateGrade(
	ctx context.Context,
	creatorID, studentID uuid.UUID,
	req *grade.CreateGradeRequest,
) (*grade.GradeResponse, error) {
	// Проверяем права создателя
	creator, err := s.userRepo.GetByID(ctx, creatorID.String())
	if err != nil {
		return nil, err
	}
	if creator == nil {
		return nil, errors.New("creator not found")
	}

	// Только преподаватель, держатель или админ могут создавать оценки
	if !creator.IsAdmin() && !creator.IsTeacher() && !creator.IsHolder() {
		return nil, errors.New("insufficient permissions to create grades")
	}

	// Проверяем существование студента
	student, err := s.userRepo.GetByID(ctx, studentID.String())
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("student not found")
	}
	if !student.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// Проверяем существование курса
	c, err := s.courseRepo.GetByID(ctx, req.CourseID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("course not found")
	}

	// Валидация оценки
	if req.Value > req.MaxValue {
		return nil, errors.New("value cannot exceed max_value")
	}

	// Создаём оценку
	g := &grade.Grade{
		UserID:     studentID,
		CourseID:   req.CourseID,
		Type:       req.Type,
		Value:      req.Value,
		MaxValue:   req.MaxValue,
		Weight:     req.Weight,
		Comment:    req.Comment,
		Date:       req.Date,
		SourceType: "manual",
		CreatedBy:  creatorID,
	}

	if err := s.gradeRepo.Create(ctx, g); err != nil {
		return nil, err
	}

	return g.ToResponse(c.Name), nil
}

// UpdateGrade обновляет оценку
func (s *GradeService) UpdateGrade(
	ctx context.Context,
	userID, gradeID uuid.UUID,
	req *grade.UpdateGradeRequest,
) (*grade.GradeResponse, error) {
	// Получаем оценку
	g, err := s.gradeRepo.GetByID(ctx, gradeID)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.New("grade not found")
	}

	// Проверяем права (только создатель или админ)
	if !s.canModifyGrade(ctx, userID, g) {
		return nil, errors.New("insufficient permissions")
	}

	// Обновляем поля
	if req.Value != nil {
		g.Value = *req.Value
	}
	if req.Comment != nil {
		g.Comment = *req.Comment
	}
	if req.Date != nil {
		g.Date = *req.Date
	}

	// Сохраняем
	if err := s.gradeRepo.Update(ctx, g); err != nil {
		return nil, err
	}

	// Получаем название курса для ответа
	c, _ := s.courseRepo.GetByID(ctx, g.CourseID)
	courseName := "Неизвестный курс"
	if c != nil {
		courseName = c.Name
	}

	return g.ToResponse(courseName), nil
}

// DeleteGrade удаляет оценку
func (s *GradeService) DeleteGrade(ctx context.Context, userID, gradeID uuid.UUID) error {
	// Получаем оценку
	g, err := s.gradeRepo.GetByID(ctx, gradeID)
	if err != nil {
		return err
	}
	if g == nil {
		return errors.New("grade not found")
	}

	// Проверяем права
	if !s.canModifyGrade(ctx, userID, g) {
		return errors.New("insufficient permissions")
	}

	return s.gradeRepo.Delete(ctx, gradeID)
}

// canModifyGrade проверяет, может ли пользователь изменять/удалять оценку
func (s *GradeService) canModifyGrade(ctx context.Context, userID uuid.UUID, g *grade.Grade) bool {
	// Получаем пользователя
	u, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil || u == nil {
		return false
	}

	// Админ может всё
	if u.IsAdmin() {
		return true
	}

	// Создатель может изменять свою оценку
	if g.CreatedBy == userID {
		return true
	}

	// Преподаватель/держатель может изменять оценки своих студентов
	// Здесь нужна дополнительная логика проверки связей
	// Пока упростим - только создатель и админ

	return false
}
