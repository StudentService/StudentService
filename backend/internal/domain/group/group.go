package group

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	CourseID   uuid.UUID  `json:"course_id"`
	SemesterID uuid.UUID  `json:"semester_id"`
	HolderID   *uuid.UUID `json:"holder_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Repository interface for group
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Group, error)
	GetByStudent(ctx context.Context, studentID uuid.UUID) (*Group, error)
	GetAll(ctx context.Context) ([]*Group, error)
	Create(ctx context.Context, group *Group) error
	Update(ctx context.Context, group *Group) error
	Delete(ctx context.Context, id uuid.UUID) error
}
