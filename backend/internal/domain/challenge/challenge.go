package challenge

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusOverdue   Status = "overdue"
	StatusDraft     Status = "draft"
)

type Challenge struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Goal        string    `json:"goal"` // цель вызова

	// Временные метки
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	Status Status `json:"status"`

	// Прогресс
	Progress int `json:"progress"` // 0-100%

	// Чекпоинты будут в отдельной таблице
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Checkpoint represents a milestone in the challenge
type Checkpoint struct {
	ID          uuid.UUID  `json:"id"`
	ChallengeID uuid.UUID  `json:"challenge_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     time.Time  `json:"due_date"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	OrderNum    int        `json:"order_num"` // для сортировки
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Artifact represents a file or link attached to challenge
type Artifact struct {
	ID          uuid.UUID `json:"id"`
	ChallengeID uuid.UUID `json:"challenge_id"`
	Type        string    `json:"type"` // "file" или "link"
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

// SelfAssessment самооценка студента
type SelfAssessment struct {
	ID          uuid.UUID `json:"id"`
	ChallengeID uuid.UUID `json:"challenge_id"`
	Rating      int       `json:"rating"` // 1-5 или 0-10
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository interface {
	// Для студента
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Challenge, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Challenge, error)

	// CRUD для вызовов
	Create(ctx context.Context, challenge *Challenge) error
	Update(ctx context.Context, challenge *Challenge) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Чекпоинты
	GetCheckpoints(ctx context.Context, challengeID uuid.UUID) ([]*Checkpoint, error)
	CreateCheckpoint(ctx context.Context, checkpoint *Checkpoint) error
	UpdateCheckpoint(ctx context.Context, checkpoint *Checkpoint) error
	DeleteCheckpoint(ctx context.Context, id uuid.UUID) error

	// Артефакты
	GetArtifacts(ctx context.Context, challengeID uuid.UUID) ([]*Artifact, error)
	CreateArtifact(ctx context.Context, artifact *Artifact) error
	DeleteArtifact(ctx context.Context, id uuid.UUID) error

	// Самооценка
	GetSelfAssessment(ctx context.Context, challengeID uuid.UUID) (*SelfAssessment, error)
	CreateSelfAssessment(ctx context.Context, assessment *SelfAssessment) error
	UpdateSelfAssessment(ctx context.Context, assessment *SelfAssessment) error
}
