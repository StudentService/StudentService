package application

import (
	"context"
	"errors"
	"time"

	"backend/internal/domain/activity"
	"backend/internal/domain/course"
	"backend/internal/domain/group"
	"backend/internal/domain/user"

	"github.com/google/uuid"
)

type ActivityService struct {
	activityRepo activity.Repository
	userRepo     user.Repository
	courseRepo   course.Repository
	groupRepo    group.Repository
}

func NewActivityService(
	activityRepo activity.Repository,
	userRepo user.Repository,
	courseRepo course.Repository,
	groupRepo group.Repository,
) *ActivityService {
	return &ActivityService{
		activityRepo: activityRepo,
		userRepo:     userRepo,
		courseRepo:   courseRepo,
		groupRepo:    groupRepo,
	}
}

// GetAvailableActivities получает активности, доступные для текущего пользователя
func (s *ActivityService) GetAvailableActivities(ctx context.Context, userID uuid.UUID) ([]*activity.ActivityResponse, error) {
	// Получаем пользователя
	u, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Получаем доступные активности
	activities, err := s.activityRepo.GetAvailableActivities(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Для каждой активности проверяем, записан ли пользователь
	responses := make([]*activity.ActivityResponse, len(activities))
	for i, a := range activities {
		participation, _ := s.activityRepo.GetParticipationByActivity(ctx, userID, a.ID)
		isEnrolled := participation != nil
		status := activity.ParticipationStatus("")
		if isEnrolled {
			status = participation.Status
		}
		responses[i] = a.ToResponse(isEnrolled, status)
	}

	return responses, nil
}

// GetMyParticipations получает все участия текущего пользователя
func (s *ActivityService) GetMyParticipations(ctx context.Context, userID uuid.UUID) ([]*activity.ParticipationResponse, error) {
	// Получаем все участия
	participations, err := s.activityRepo.GetMyParticipations(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Для каждого участия получаем детали активности
	responses := make([]*activity.ParticipationResponse, len(participations))
	for i, p := range participations {
		a, err := s.activityRepo.GetByID(ctx, p.ActivityID)
		if err != nil {
			continue
		}
		if a == nil {
			continue
		}
		responses[i] = p.ToResponse(a)
	}

	return responses, nil
}

// Enroll записывает текущего пользователя на активность
func (s *ActivityService) Enroll(ctx context.Context, userID uuid.UUID, activityID uuid.UUID) error {
	// Проверяем существование активности
	a, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("activity not found")
	}

	// Проверяем, активна ли активность
	if a.Status != activity.ActivityStatusActive {
		return errors.New("activity is not active")
	}

	// Проверяем, не записан ли уже пользователь
	existing, err := s.activityRepo.GetParticipationByActivity(ctx, userID, activityID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("already enrolled in this activity")
	}

	// Проверяем, есть ли места
	if a.MaxParticipants > 0 && a.CurrentParticipants >= a.MaxParticipants {
		return errors.New("no available slots")
	}

	// Создаём участие
	participation := &activity.Participation{
		ActivityID:   activityID,
		UserID:       userID,
		Status:       activity.ParticipationStatusEnrolled,
		PointsEarned: 0,
	}

	return s.activityRepo.Enroll(ctx, participation)
}

// CancelEnrollment отменяет запись на активность
func (s *ActivityService) CancelEnrollment(ctx context.Context, userID uuid.UUID, activityID uuid.UUID) error {
	// Получаем участие
	p, err := s.activityRepo.GetParticipationByActivity(ctx, userID, activityID)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("not enrolled in this activity")
	}

	// Проверяем, можно ли отменить (не началась ли активность)
	a, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return err
	}
	if a.StartTime != nil && a.StartTime.Before(time.Now()) {
		return errors.New("cannot cancel enrollment after activity has started")
	}

	return s.activityRepo.CancelEnrollment(ctx, p.ID)
}

// CreateActivity создаёт новую активность (для преподавателя/админа)
func (s *ActivityService) CreateActivity(ctx context.Context, creatorID uuid.UUID, req *activity.CreateActivityRequest) (*activity.Activity, error) {
	// Получаем создателя
	creator, err := s.userRepo.GetByID(ctx, creatorID.String())
	if err != nil {
		return nil, err
	}
	if creator == nil {
		return nil, errors.New("creator not found")
	}

	// Проверяем права
	if !s.canManageActivities(creator) {
		return nil, errors.New("insufficient permissions to create activities")
	}

	// Валидация дат
	if req.StartTime != nil && req.EndTime != nil && req.StartTime.After(*req.EndTime) {
		return nil, errors.New("start time must be before end time")
	}

	// Проверяем существование курса/группы, если указаны
	if req.CourseID != nil {
		course, err := s.courseRepo.GetByID(ctx, *req.CourseID)
		if err != nil || course == nil {
			return nil, errors.New("course not found")
		}
	}
	if req.GroupID != nil {
		group, err := s.groupRepo.GetByID(ctx, *req.GroupID)
		if err != nil || group == nil {
			return nil, errors.New("group not found")
		}
	}

	// Создаём активность
	a := &activity.Activity{
		Title:           req.Title,
		Description:     req.Description,
		Type:            req.Type,
		Status:          activity.ActivityStatusActive,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Deadline:        req.Deadline,
		Location:        req.Location,
		OnlineLink:      req.OnlineLink,
		MaxParticipants: req.MaxParticipants,
		Points:          req.Points,
		Weight:          req.Weight,
		CourseID:        req.CourseID,
		GroupID:         req.GroupID,
		CreatedBy:       creatorID,
		CreatedByRole:   string(creator.Role),
	}

	if err := s.activityRepo.Create(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

// UpdateActivity обновляет активность
func (s *ActivityService) UpdateActivity(ctx context.Context, userID uuid.UUID, activityID uuid.UUID, req *activity.UpdateActivityRequest) (*activity.Activity, error) {
	// Получаем активность
	a, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("activity not found")
	}

	// Проверяем права
	if !s.canModifyActivity(ctx, userID, a) {
		return nil, errors.New("insufficient permissions")
	}

	// Обновляем поля
	if req.Title != nil {
		a.Title = *req.Title
	}
	if req.Description != nil {
		a.Description = *req.Description
	}
	if req.Type != nil {
		a.Type = *req.Type
	}
	if req.Status != nil {
		a.Status = *req.Status
	}
	if req.StartTime != nil {
		a.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		a.EndTime = req.EndTime
	}
	if req.Deadline != nil {
		a.Deadline = req.Deadline
	}
	if req.Location != nil {
		a.Location = *req.Location
	}
	if req.OnlineLink != nil {
		a.OnlineLink = *req.OnlineLink
	}
	if req.MaxParticipants != nil {
		a.MaxParticipants = *req.MaxParticipants
	}
	if req.Points != nil {
		a.Points = *req.Points
	}
	if req.Weight != nil {
		a.Weight = *req.Weight
	}

	if err := s.activityRepo.Update(ctx, a); err != nil {
		return nil, err
	}

	return a, nil
}

// DeleteActivity удаляет активность
func (s *ActivityService) DeleteActivity(ctx context.Context, userID uuid.UUID, activityID uuid.UUID) error {
	// Получаем активность
	a, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("activity not found")
	}

	// Проверяем права
	if !s.canModifyActivity(ctx, userID, a) {
		return errors.New("insufficient permissions")
	}

	// Проверяем, есть ли участники
	if a.CurrentParticipants > 0 {
		return errors.New("cannot delete activity with enrolled participants")
	}

	return s.activityRepo.Delete(ctx, activityID)
}

// GetActivityParticipants получает список участников активности (для преподавателя)
func (s *ActivityService) GetActivityParticipants(ctx context.Context, userID uuid.UUID, activityID uuid.UUID) ([]*activity.ParticipationResponse, error) {
	// Получаем активность
	a, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("activity not found")
	}

	// Проверяем права
	if !s.canViewParticipants(ctx, userID, a) {
		return nil, errors.New("insufficient permissions")
	}

	// Получаем участников
	participations, err := s.activityRepo.GetActivityParticipants(ctx, activityID)
	if err != nil {
		return nil, err
	}

	responses := make([]*activity.ParticipationResponse, len(participations))
	for i, p := range participations {
		responses[i] = p.ToResponse(a)
	}

	return responses, nil
}

// MarkAttendance отмечает посещение
func (s *ActivityService) MarkAttendance(ctx context.Context, userID uuid.UUID, participationID uuid.UUID, attended bool) error {
	// Получаем участие
	// Для этого нужен метод GetParticipationByID, добавим позже при необходимости
	// Пока используем прямой вызов репозитория с проверкой через активность
	return errors.New("not implemented")
}

// SetGrade выставляет оценку за активность
func (s *ActivityService) SetGrade(ctx context.Context, userID uuid.UUID, participationID uuid.UUID, grade float64, feedback string) error {
	// TODO: добавить проверку прав
	return s.activityRepo.SetGrade(ctx, participationID, grade, feedback)
}

// canManageActivities проверяет, может ли пользователь создавать активности
func (s *ActivityService) canManageActivities(u *user.User) bool {
	return u.IsAdmin() || u.IsTeacher() || u.IsHolder()
}

// canModifyActivity проверяет, может ли пользователь изменять активность
func (s *ActivityService) canModifyActivity(ctx context.Context, userID uuid.UUID, a *activity.Activity) bool {
	u, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil || u == nil {
		return false
	}

	// Админ может всё
	if u.IsAdmin() {
		return true
	}

	// Создатель может изменять свою активность
	if a.CreatedBy == userID {
		return true
	}

	// Преподаватель может изменять активности своего курса
	if u.IsTeacher() && a.CourseID != nil {
		// Здесь нужна проверка, ведёт ли преподаватель этот курс
		return true
	}

	// Держатель может изменять активности своей группы
	if u.IsHolder() && a.GroupID != nil {
		// Здесь нужна проверка, является ли держатель holder_id группы
		return true
	}

	return false
}

// canViewParticipants проверяет, может ли пользователь видеть участников
func (s *ActivityService) canViewParticipants(ctx context.Context, userID uuid.UUID, a *activity.Activity) bool {
	u, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil || u == nil {
		return false
	}

	// Админ видит всех
	if u.IsAdmin() {
		return true
	}

	// Создатель видит участников
	if a.CreatedBy == userID {
		return true
	}

	// Преподаватель видит участников своего курса
	if u.IsTeacher() && a.CourseID != nil {
		return true
	}

	// Держатель видит участников своей группы
	if u.IsHolder() && a.GroupID != nil {
		return true
	}

	return false
}
