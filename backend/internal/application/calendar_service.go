package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend/internal/domain/calendar"
	"backend/internal/domain/course"
	"backend/internal/domain/group"
	"backend/internal/domain/user"
)

type CalendarService struct {
	eventRepo  calendar.Repository
	userRepo   user.Repository
	courseRepo course.Repository
	groupRepo  group.Repository
}

func NewCalendarService(
	eventRepo calendar.Repository,
	userRepo user.Repository,
	courseRepo course.Repository,
	groupRepo group.Repository,
) *CalendarService {
	return &CalendarService{
		eventRepo:  eventRepo,
		userRepo:   userRepo,
		courseRepo: courseRepo,
		groupRepo:  groupRepo,
	}
}

// GetMyEvents получает события текущего пользователя за период
func (s *CalendarService) GetMyEvents(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*calendar.EventResponse, error) {
	// Получаем пользователя
	u, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Получаем события
	events, err := s.eventRepo.GetStudentEvents(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	// Конвертируем в Response DTO, подгружая дополнительные данные
	responses := make([]*calendar.EventResponse, len(events))
	for i, e := range events {
		// Получаем названия курса и группы
		courseName, groupName, creatorName := s.getEventNames(ctx, e, nil)

		responses[i] = &calendar.EventResponse{
			ID:            e.ID,
			Title:         e.Title,
			Description:   e.Description,
			Type:          e.Type,
			StartTime:     e.StartTime,
			EndTime:       e.EndTime,
			AllDay:        e.AllDay,
			Location:      e.Location,
			OnlineLink:    e.OnlineLink,
			CourseName:    courseName,
			GroupName:     groupName,
			CreatedBy:     creatorName,
			CreatedByRole: e.CreatedByRole,
		}
	}

	return responses, nil
}

// CreateEvent создаёт новое событие
func (s *CalendarService) CreateEvent(ctx context.Context, creatorID uuid.UUID, req *calendar.CreateEventRequest) (*calendar.EventResponse, error) {
	// Получаем создателя
	creator, err := s.userRepo.GetByID(ctx, creatorID.String())
	if err != nil {
		return nil, err
	}
	if creator == nil {
		return nil, errors.New("creator not found")
	}

	// Проверяем права на создание событий
	if !s.canCreateEvents(creator) {
		return nil, errors.New("insufficient permissions to create events")
	}

	// Валидация - событие должно быть для кого-то
	if req.CourseID == nil && req.GroupID == nil && req.UserID == nil {
		return nil, errors.New("event must be assigned to course, group or user")
	}

	// Если событие для конкретного пользователя, проверяем что он студент
	if req.UserID != nil {
		student, err := s.userRepo.GetByID(ctx, req.UserID.String())
		if err != nil {
			return nil, err
		}
		if student == nil {
			return nil, errors.New("student not found")
		}
		if !student.IsStudent() {
			return nil, errors.New("user is not a student")
		}
	}

	// Создаём событие
	event := &calendar.Event{
		Title:         req.Title,
		Description:   req.Description,
		Type:          req.Type,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		AllDay:        req.AllDay,
		Location:      req.Location,
		OnlineLink:    req.OnlineLink,
		CourseID:      req.CourseID,
		GroupID:       req.GroupID,
		UserID:        req.UserID,
		CreatedBy:     creatorID,
		CreatedByRole: string(creator.Role),
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	// Получаем названия для ответа
	courseName, groupName, creatorName := s.getEventNames(ctx, event, creator)

	return &calendar.EventResponse{
		ID:            event.ID,
		Title:         event.Title,
		Description:   event.Description,
		Type:          event.Type,
		StartTime:     event.StartTime,
		EndTime:       event.EndTime,
		AllDay:        event.AllDay,
		Location:      event.Location,
		OnlineLink:    event.OnlineLink,
		CourseName:    courseName,
		GroupName:     groupName,
		CreatedBy:     creatorName,
		CreatedByRole: event.CreatedByRole,
	}, nil
}

// UpdateEvent обновляет событие
func (s *CalendarService) UpdateEvent(ctx context.Context, userID, eventID uuid.UUID, req *calendar.UpdateEventRequest) (*calendar.EventResponse, error) {
	// Получаем событие
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event not found")
	}

	// Проверяем права на редактирование
	if !s.canModifyEvent(ctx, userID, event) {
		return nil, errors.New("insufficient permissions to modify this event")
	}

	// Обновляем поля
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Type != nil {
		event.Type = *req.Type
	}
	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		event.EndTime = *req.EndTime
	}
	if req.AllDay != nil {
		event.AllDay = *req.AllDay
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.OnlineLink != nil {
		event.OnlineLink = *req.OnlineLink
	}
	if req.CourseID != nil {
		event.CourseID = req.CourseID
	}
	if req.GroupID != nil {
		event.GroupID = req.GroupID
	}
	if req.UserID != nil {
		event.UserID = req.UserID
	}

	// Сохраняем
	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	// Получаем названия для ответа
	courseName, groupName, creatorName := s.getEventNames(ctx, event, nil)

	return &calendar.EventResponse{
		ID:            event.ID,
		Title:         event.Title,
		Description:   event.Description,
		Type:          event.Type,
		StartTime:     event.StartTime,
		EndTime:       event.EndTime,
		AllDay:        event.AllDay,
		Location:      event.Location,
		OnlineLink:    event.OnlineLink,
		CourseName:    courseName,
		GroupName:     groupName,
		CreatedBy:     creatorName,
		CreatedByRole: event.CreatedByRole,
	}, nil
}

// DeleteEvent удаляет событие
func (s *CalendarService) DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	// Получаем событие
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}

	// Проверяем права на удаление
	if !s.canModifyEvent(ctx, userID, event) {
		return errors.New("insufficient permissions to delete this event")
	}

	return s.eventRepo.Delete(ctx, eventID)
}

// canCreateEvents проверяет, может ли пользователь создавать события
func (s *CalendarService) canCreateEvents(u *user.User) bool {
	return u.IsAdmin() || u.IsTeacher() || u.IsHolder()
}

// canModifyEvent проверяет, может ли пользователь изменять/удалять событие
func (s *CalendarService) canModifyEvent(ctx context.Context, userID uuid.UUID, event *calendar.Event) bool {
	// Получаем пользователя
	u, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil || u == nil {
		return false
	}

	// Админ может всё
	if u.IsAdmin() {
		return true
	}

	// Создатель может изменять своё событие
	if event.CreatedBy == userID {
		return true
	}

	return false
}

// getEventNames получает названия связанных сущностей для ответа
func (s *CalendarService) getEventNames(ctx context.Context, event *calendar.Event, creator *user.User) (courseName, groupName, creatorName string) {
	if event.CourseID != nil {
		if course, err := s.courseRepo.GetByID(ctx, *event.CourseID); err == nil && course != nil {
			courseName = course.Name
		}
	}

	if event.GroupID != nil {
		if group, err := s.groupRepo.GetByID(ctx, *event.GroupID); err == nil && group != nil {
			groupName = group.Name
		}
	}

	if creator != nil {
		creatorName = creator.GetFullName()
	} else if event.CreatedBy != uuid.Nil {
		if c, err := s.userRepo.GetByID(ctx, event.CreatedBy.String()); err == nil && c != nil {
			creatorName = c.GetFullName()
		}
	}

	return
}
