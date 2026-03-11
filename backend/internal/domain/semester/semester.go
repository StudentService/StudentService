package semester

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Semester struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"` // например "Осень 2026"
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository interface for semester
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Semester, error)
	GetActive(ctx context.Context) (*Semester, error)
	GetAll(ctx context.Context) ([]*Semester, error)
	Create(ctx context.Context, semester *Semester) error
	Update(ctx context.Context, semester *Semester) error
	Delete(ctx context.Context, id uuid.UUID) error
}
