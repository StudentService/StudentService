package calendar

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeClass    EventType = "class"    // занятие
	EventTypeMeeting  EventType = "meeting"  // встреча
	EventTypeDeadline EventType = "deadline" // дедлайн
	EventTypeActivity EventType = "activity" // мероприятие
	EventTypeExam     EventType = "exam"     // экзамен
	EventTypeHoliday  EventType = "holiday"  // выходной/праздник
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        EventType `json:"type"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	AllDay      bool      `json:"all_day"`
	Location    string    `json:"location,omitempty"`
	OnlineLink  string    `json:"online_link,omitempty"`

	// Связи
	CourseID      *uuid.UUID `json:"course_id,omitempty"`
	GroupID       *uuid.UUID `json:"group_id,omitempty"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedByRole string     `json:"created_by_role"` // кто создал (admin/teacher/holder)

	// Для личных событий
	UserID *uuid.UUID `json:"user_id,omitempty"` // если событие привязано к конкретному студенту

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EventResponse DTO для ответа
type EventResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Type          EventType `json:"type"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	AllDay        bool      `json:"all_day"`
	Location      string    `json:"location,omitempty"`
	OnlineLink    string    `json:"online_link,omitempty"`
	CourseName    string    `json:"course_name,omitempty"`
	GroupName     string    `json:"group_name,omitempty"`
	CreatedBy     string    `json:"created_by"` // имя создателя
	CreatedByRole string    `json:"created_by_role"`
}

// ToResponse конвертирует Event в EventResponse
func (e *Event) ToResponse(courseName, groupName, creatorName string) *EventResponse {
	return &EventResponse{
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

type Repository interface {
	// Для студентов - получить события, которые им доступны
	GetStudentEvents(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*Event, error)

	// Для групп - получить события группы
	GetGroupEvents(ctx context.Context, groupID uuid.UUID, from, to time.Time) ([]*Event, error)

	// Для курсов - получить события курса
	GetCourseEvents(ctx context.Context, courseID uuid.UUID, from, to time.Time) ([]*Event, error)

	// CRUD для событий (админ/преподаватель)
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	Create(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}
