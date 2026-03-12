package activity

import (
	"time"

	"github.com/google/uuid"
)

type ActivityType string

const (
	ActivityTypeClass    ActivityType = "class"    // занятие
	ActivityTypeWorkshop ActivityType = "workshop" // воркшоп
	ActivityTypeMeeting  ActivityType = "meeting"  // встреча
	ActivityTypeTask     ActivityType = "task"     // задача
	ActivityTypeProject  ActivityType = "project"  // проект
	ActivityTypeEvent    ActivityType = "event"    // мероприятие
)

type ActivityStatus string

const (
	ActivityStatusActive    ActivityStatus = "active"    // активна
	ActivityStatusCompleted ActivityStatus = "completed" // завершена
	ActivityStatusCancelled ActivityStatus = "cancelled" // отменена
	ActivityStatusDraft     ActivityStatus = "draft"     // черновик
)

// Activity - активность в каталоге
type Activity struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        ActivityType   `json:"type"`
	Status      ActivityStatus `json:"status"`

	// Временные параметры
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Deadline  *time.Time `json:"deadline,omitempty"` // для задач

	// Место проведения
	Location   string `json:"location,omitempty"`
	OnlineLink string `json:"online_link,omitempty"`

	// Ограничения
	MaxParticipants     int `json:"max_participants,omitempty"`
	CurrentParticipants int `json:"current_participants,omitempty"`

	// Баллы/вес
	Points int     `json:"points"` // баллы за участие
	Weight float64 `json:"weight"` // вес для метрик

	// Для кого активность
	CourseID      *uuid.UUID `json:"course_id,omitempty"`
	GroupID       *uuid.UUID `json:"group_id,omitempty"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedByRole string     `json:"created_by_role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Participation - участие студента в активности
type Participation struct {
	ID         uuid.UUID           `json:"id"`
	ActivityID uuid.UUID           `json:"activity_id"`
	UserID     uuid.UUID           `json:"user_id"`
	Status     ParticipationStatus `json:"status"`

	// Результаты
	Grade        *float64 `json:"grade,omitempty"`    // оценка
	Feedback     string   `json:"feedback,omitempty"` // обратная связь
	PointsEarned int      `json:"points_earned"`      // полученные баллы

	// Временные метки
	EnrolledAt  time.Time  `json:"enrolled_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ParticipationStatus string

const (
	ParticipationStatusEnrolled  ParticipationStatus = "enrolled"  // записан
	ParticipationStatusAttended  ParticipationStatus = "attended"  // посетил
	ParticipationStatusCompleted ParticipationStatus = "completed" // выполнил
	ParticipationStatusMissed    ParticipationStatus = "missed"    // пропустил
	ParticipationStatusCancelled ParticipationStatus = "cancelled" // отменил участие
)
