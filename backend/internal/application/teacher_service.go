package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend/internal/domain/activity"
	"backend/internal/domain/grade"
	"backend/internal/domain/group"
	"backend/internal/domain/teacher"
	"backend/internal/domain/user"
)

type TeacherService struct {
	userRepo     user.Repository
	groupRepo    group.Repository
	gradeRepo    grade.Repository
	activityRepo activity.Repository
}

func NewTeacherService(
	userRepo user.Repository,
	groupRepo group.Repository,
	gradeRepo grade.Repository,
	activityRepo activity.Repository,
) *TeacherService {
	return &TeacherService{
		userRepo:     userRepo,
		groupRepo:    groupRepo,
		gradeRepo:    gradeRepo,
		activityRepo: activityRepo,
	}
}

// GetDashboard получает дашборд преподавателя
func (s *TeacherService) GetDashboard(ctx context.Context, teacherID uuid.UUID) (*teacher.DashboardResponse, error) {
	// Получаем группы преподавателя
	groups, err := s.getTeacherGroups(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// Получаем недавние активности
	recentActivities, err := s.getRecentActivities(ctx, teacherID)
	if err != nil {
		recentActivities = []*teacher.ActivitySummary{}
	}

	// Подсчитываем общее количество студентов
	totalStudents := 0
	for _, g := range groups {
		totalStudents += g.StudentsCount
	}

	return &teacher.DashboardResponse{
		Groups:           groups,
		RecentActivities: recentActivities,
		PendingReviews:   0, // TODO: реализовать подсчёт
		TotalStudents:    totalStudents,
	}, nil
}

// GetTeacherGroups получает группы преподавателя
func (s *TeacherService) GetTeacherGroups(ctx context.Context, teacherID uuid.UUID) ([]*teacher.GroupSummary, error) {
	return s.getTeacherGroups(ctx, teacherID)
}

// GetGroupStudents получает студентов группы
func (s *TeacherService) GetGroupStudents(ctx context.Context, teacherID, groupID uuid.UUID) ([]*user.UserResponse, error) {
	// Проверяем, что преподаватель ведёт эту группу
	if !s.isTeacherOfGroup(ctx, teacherID, groupID) {
		return nil, errors.New("access denied")
	}

	students, err := s.userRepo.GetByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	responses := make([]*user.UserResponse, len(students))
	for i, u := range students {
		responses[i] = u.ToResponse()
	}
	return responses, nil
}

// GetStudentProfile получает профиль студента
func (s *TeacherService) GetStudentProfile(ctx context.Context, teacherID, studentID uuid.UUID) (*teacher.StudentProfile, error) {
	// Получаем студента
	student, err := s.userRepo.GetByID(ctx, studentID.String())
	if err != nil {
		return nil, err
	}
	if student == nil || !student.IsStudent() {
		return nil, errors.New("student not found")
	}

	// Проверяем, что преподаватель имеет доступ к этому студенту
	if !s.canAccessStudent(ctx, teacherID, studentID) {
		return nil, errors.New("access denied")
	}

	// Получаем среднюю оценку
	avgGrade, _ := s.gradeRepo.GetAverageByUser(ctx, studentID)

	// Получаем группу
	groupName := ""
	if student.GroupID != nil {
		g, err := s.groupRepo.GetByID(ctx, *student.GroupID)
		if err == nil && g != nil {
			groupName = g.Name
		}
	}

	return &teacher.StudentProfile{
		ID:           student.ID,
		FirstName:    student.FirstName,
		LastName:     student.LastName,
		Email:        student.Email,
		GroupName:    groupName,
		AverageGrade: avgGrade,
		TotalPoints:  0, // TODO: рассчитать из активностей
	}, nil
}

// GetStudentGrades получает оценки студента
func (s *TeacherService) GetStudentGrades(ctx context.Context, teacherID, studentID uuid.UUID) ([]*teacher.GradeWithStudent, error) {
	// Проверяем доступ
	if !s.canAccessStudent(ctx, teacherID, studentID) {
		return nil, errors.New("access denied")
	}

	grades, err := s.gradeRepo.GetByUserID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	student, _ := s.userRepo.GetByID(ctx, studentID.String())
	studentName := ""
	if student != nil {
		studentName = student.GetFullName()
	}

	result := make([]*teacher.GradeWithStudent, len(grades))
	for i, g := range grades {
		result[i] = &teacher.GradeWithStudent{
			ID:          g.ID,
			StudentID:   g.UserID,
			StudentName: studentName,
			Value:       g.Value,
			MaxValue:    g.MaxValue,
			Type:        string(g.Type),
			Date:        g.Date,
			Comment:     g.Comment,
		}
	}
	return result, nil
}

// GetStudentChallenges получает вызовы студента
func (s *TeacherService) GetStudentChallenges(ctx context.Context, teacherID, studentID uuid.UUID) (interface{}, error) {
	// Проверяем доступ
	if !s.canAccessStudent(ctx, teacherID, studentID) {
		return nil, errors.New("access denied")
	}

	// TODO: реализовать получение вызовов студента
	return nil, nil
}

// GetTeacherActivities получает активности преподавателя
func (s *TeacherService) GetTeacherActivities(ctx context.Context, teacherID uuid.UUID) ([]*activity.ActivityResponse, error) {
	activities, err := s.activityRepo.GetByCreator(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	responses := make([]*activity.ActivityResponse, len(activities))
	for i, a := range activities {
		responses[i] = a.ToResponse(false, "")
	}
	return responses, nil
}

// GetGroupGrades получает оценки группы
func (s *TeacherService) GetGroupGrades(ctx context.Context, teacherID, groupID uuid.UUID) ([]*teacher.GradeWithStudent, error) {
	// Проверяем доступ к группе
	if !s.isTeacherOfGroup(ctx, teacherID, groupID) {
		return nil, errors.New("access denied")
	}

	// Получаем студентов группы
	students, err := s.userRepo.GetByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var allGrades []*teacher.GradeWithStudent
	for _, student := range students {
		grades, _ := s.gradeRepo.GetByUserID(ctx, student.ID)
		for _, g := range grades {
			allGrades = append(allGrades, &teacher.GradeWithStudent{
				ID:          g.ID,
				StudentID:   g.UserID,
				StudentName: student.GetFullName(),
				Value:       g.Value,
				MaxValue:    g.MaxValue,
				Type:        string(g.Type),
				Date:        g.Date,
				Comment:     g.Comment,
			})
		}
	}
	return allGrades, nil
}

// ImportGrades импортирует оценки
func (s *TeacherService) ImportGrades(ctx context.Context, teacherID uuid.UUID, req *teacher.ImportGradesRequest) error {
	// Проверяем доступ к группе
	if !s.isTeacherOfGroup(ctx, teacherID, req.GroupID) {
		return errors.New("access denied")
	}

	// Для каждого студента создаём оценку
	for _, item := range req.Grades {
		// Проверяем, что студент в этой группе
		student, err := s.userRepo.GetByID(ctx, item.StudentID.String())
		if err != nil || student == nil || student.GroupID == nil || *student.GroupID != req.GroupID {
			continue
		}

		g := &grade.Grade{
			UserID:     item.StudentID,
			CourseID:   uuid.Nil,            // TODO: определить по группе
			Type:       grade.GradeTypeExam, // TODO: определить из запроса
			Value:      item.Value,
			MaxValue:   item.MaxValue,
			Comment:    item.Comment,
			Date:       time.Now(),
			SourceType: "import",
			CreatedBy:  teacherID,
		}
		s.gradeRepo.Create(ctx, g)
	}
	return nil
}

// MarkAttendance отмечает посещаемость
func (s *TeacherService) MarkAttendance(ctx context.Context, teacherID, activityID uuid.UUID, req *teacher.AttendanceRequest) error {
	// Проверяем, что преподаватель создал эту активность
	act, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return err
	}
	if act == nil {
		return errors.New("activity not found")
	}
	if act.CreatedBy != teacherID {
		return errors.New("access denied")
	}

	// Отмечаем каждого студента
	for _, studentID := range req.StudentIDs {
		participation, err := s.activityRepo.GetParticipationByActivity(ctx, studentID, activityID)
		if err != nil {
			continue // пропускаем ошибки
		}
		if participation != nil {
			var status activity.ParticipationStatus
			if req.Attended {
				status = activity.ParticipationStatusAttended
			} else {
				status = activity.ParticipationStatusMissed
			}
			// Используем существующий метод UpdateParticipationStatus
			s.activityRepo.UpdateParticipationStatus(ctx, participation.ID, status, nil, "")
		}
	}
	return nil
}

// Вспомогательные методы

func (s *TeacherService) getTeacherGroups(ctx context.Context, teacherID uuid.UUID) ([]*teacher.GroupSummary, error) {
	// TODO: получить группы, где преподаватель ведёт
	// Пока возвращаем пустой список
	return []*teacher.GroupSummary{}, nil
}

func (s *TeacherService) getRecentActivities(ctx context.Context, teacherID uuid.UUID) ([]*teacher.ActivitySummary, error) {
	activities, err := s.activityRepo.GetByCreator(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	var summaries []*teacher.ActivitySummary
	for _, a := range activities {
		summaries = append(summaries, &teacher.ActivitySummary{
			ID:           a.ID,
			Title:        a.Title,
			Type:         string(a.Type),
			Date:         *a.StartTime,
			Participants: a.CurrentParticipants,
		})
	}
	return summaries, nil
}

func (s *TeacherService) isTeacherOfGroup(ctx context.Context, teacherID, groupID uuid.UUID) bool {
	// TODO: реализовать проверку
	return true // временно
}

func (s *TeacherService) canAccessStudent(ctx context.Context, teacherID, studentID uuid.UUID) bool {
	// TODO: реализовать проверку
	return true // временно
}
