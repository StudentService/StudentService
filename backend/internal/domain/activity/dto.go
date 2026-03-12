package activity

import (
	"time"

	"github.com/google/uuid"
)

// ActivityResponse DTO для отображения активности
type ActivityResponse struct {
	ID                  uuid.UUID           `json:"id"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	Type                ActivityType        `json:"type"`
	Status              ActivityStatus      `json:"status"`
	StartTime           *time.Time          `json:"start_time,omitempty"`
	EndTime             *time.Time          `json:"end_time,omitempty"`
	Deadline            *time.Time          `json:"deadline,omitempty"`
	Location            string              `json:"location,omitempty"`
	OnlineLink          string              `json:"online_link,omitempty"`
	MaxParticipants     int                 `json:"max_participants,omitempty"`
	CurrentParticipants int                 `json:"current_participants"`
	Points              int                 `json:"points"`
	IsEnrolled          bool                `json:"is_enrolled"` // флаг, записан ли текущий пользователь
	EnrollmentStatus    ParticipationStatus `json:"enrollment_status,omitempty"`
}

// ParticipationResponse DTO для отображения участия
type ParticipationResponse struct {
	ID            uuid.UUID           `json:"id"`
	ActivityID    uuid.UUID           `json:"activity_id"`
	ActivityTitle string              `json:"activity_title"`
	ActivityType  ActivityType        `json:"activity_type"`
	Status        ParticipationStatus `json:"status"`
	Grade         *float64            `json:"grade,omitempty"`
	Feedback      string              `json:"feedback,omitempty"`
	PointsEarned  int                 `json:"points_earned"`
	EnrolledAt    time.Time           `json:"enrolled_at"`
	CompletedAt   *time.Time          `json:"completed_at,omitempty"`
	StartTime     *time.Time          `json:"start_time,omitempty"`
	Location      string              `json:"location,omitempty"`
}

// CreateActivityRequest DTO для создания активности
type CreateActivityRequest struct {
	Title           string       `json:"title" binding:"required"`
	Description     string       `json:"description"`
	Type            ActivityType `json:"type" binding:"required,oneof=class workshop meeting task project event"`
	StartTime       *time.Time   `json:"start_time,omitempty"`
	EndTime         *time.Time   `json:"end_time,omitempty"`
	Deadline        *time.Time   `json:"deadline,omitempty"`
	Location        string       `json:"location"`
	OnlineLink      string       `json:"online_link"`
	MaxParticipants int          `json:"max_participants"`
	Points          int          `json:"points"`
	Weight          float64      `json:"weight"`
	CourseID        *uuid.UUID   `json:"course_id,omitempty"`
	GroupID         *uuid.UUID   `json:"group_id,omitempty"`
}

// UpdateActivityRequest DTO для обновления активности
type UpdateActivityRequest struct {
	Title           *string         `json:"title,omitempty"`
	Description     *string         `json:"description,omitempty"`
	Type            *ActivityType   `json:"type,omitempty" binding:"omitempty,oneof=class workshop meeting task project event"`
	Status          *ActivityStatus `json:"status,omitempty" binding:"omitempty,oneof=active completed cancelled draft"`
	StartTime       *time.Time      `json:"start_time,omitempty"`
	EndTime         *time.Time      `json:"end_time,omitempty"`
	Deadline        *time.Time      `json:"deadline,omitempty"`
	Location        *string         `json:"location,omitempty"`
	OnlineLink      *string         `json:"online_link,omitempty"`
	MaxParticipants *int            `json:"max_participants,omitempty"`
	Points          *int            `json:"points,omitempty"`
	Weight          *float64        `json:"weight,omitempty"`
}

// EnrollRequest DTO для записи на активность
type EnrollRequest struct {
	ActivityID uuid.UUID `json:"activity_id" binding:"required"`
}

// GradeParticipationRequest DTO для выставления оценки
type GradeParticipationRequest struct {
	Grade    float64 `json:"grade" binding:"required"`
	Feedback string  `json:"feedback"`
}

// ToResponse конвертирует Activity в ActivityResponse
func (a *Activity) ToResponse(isEnrolled bool, status ParticipationStatus) *ActivityResponse {
	return &ActivityResponse{
		ID:                  a.ID,
		Title:               a.Title,
		Description:         a.Description,
		Type:                a.Type,
		Status:              a.Status,
		StartTime:           a.StartTime,
		EndTime:             a.EndTime,
		Deadline:            a.Deadline,
		Location:            a.Location,
		OnlineLink:          a.OnlineLink,
		MaxParticipants:     a.MaxParticipants,
		CurrentParticipants: a.CurrentParticipants,
		Points:              a.Points,
		IsEnrolled:          isEnrolled,
		EnrollmentStatus:    status,
	}
}

// ToResponse конвертирует Participation в ParticipationResponse
func (p *Participation) ToResponse(activity *Activity) *ParticipationResponse {
	return &ParticipationResponse{
		ID:            p.ID,
		ActivityID:    p.ActivityID,
		ActivityTitle: activity.Title,
		ActivityType:  activity.Type,
		Status:        p.Status,
		Grade:         p.Grade,
		Feedback:      p.Feedback,
		PointsEarned:  p.PointsEarned,
		EnrolledAt:    p.EnrolledAt,
		CompletedAt:   p.CompletedAt,
		StartTime:     activity.StartTime,
		Location:      activity.Location,
	}
}
