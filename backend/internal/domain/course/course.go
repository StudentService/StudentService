package course

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Credits     int       `json:"credits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Repository interface for course
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Course, error)
	GetAll(ctx context.Context) ([]*Course, error)
	Create(ctx context.Context, course *Course) error
	Update(ctx context.Context, course *Course) error
	Delete(ctx context.Context, id uuid.UUID) error
}
