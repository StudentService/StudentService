package calendar

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Type        EventType `json:"type" binding:"required,oneof=class meeting deadline activity exam holiday"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
	AllDay      bool      `json:"all_day"`
	Location    string    `json:"location"`
	OnlineLink  string    `json:"online_link"`

	// Для кого событие (можно указать одно из)
	CourseID *uuid.UUID `json:"course_id,omitempty"` // для всего курса
	GroupID  *uuid.UUID `json:"group_id,omitempty"`  // для группы
	UserID   *uuid.UUID `json:"user_id,omitempty"`   // для конкретного студента
}

type UpdateEventRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Type        *EventType `json:"type,omitempty" binding:"omitempty,oneof=class meeting deadline activity exam holiday"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	AllDay      *bool      `json:"all_day,omitempty"`
	Location    *string    `json:"location,omitempty"`
	OnlineLink  *string    `json:"online_link,omitempty"`
	CourseID    *uuid.UUID `json:"course_id,omitempty"`
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
}
