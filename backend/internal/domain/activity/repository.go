package activity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// Для всех пользователей
	GetAvailableActivities(ctx context.Context, userID uuid.UUID) ([]*Activity, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Activity, error)

	// Для студента - мои участия
	GetMyParticipations(ctx context.Context, userID uuid.UUID) ([]*Participation, error)
	GetParticipationByActivity(ctx context.Context, userID, activityID uuid.UUID) (*Participation, error)
	Enroll(ctx context.Context, participation *Participation) error
	CancelEnrollment(ctx context.Context, id uuid.UUID) error
	UpdateParticipationStatus(ctx context.Context, id uuid.UUID, status ParticipationStatus, grade *float64, feedback string) error

	// Для преподавателя/админа - управление активностями
	Create(ctx context.Context, activity *Activity) error
	Update(ctx context.Context, activity *Activity) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListActivities(ctx context.Context, filter ActivityFilter) ([]*Activity, error)

	// Для преподавателя - управление участиями
	GetActivityParticipants(ctx context.Context, activityID uuid.UUID) ([]*Participation, error)
	MarkAttendance(ctx context.Context, participationID uuid.UUID, attended bool) error
	SetGrade(ctx context.Context, participationID uuid.UUID, grade float64, feedback string) error
}

type ActivityFilter struct {
	Type     []ActivityType
	Status   []ActivityStatus
	CourseID *uuid.UUID
	GroupID  *uuid.UUID
	FromDate *time.Time
	ToDate   *time.Time
	Search   string
}
