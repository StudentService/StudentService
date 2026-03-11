package challenge

import (
	"time"

	"github.com/google/uuid"
)

type ChallengeResponse struct {
	ID          uuid.UUID             `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Goal        string                `json:"goal"`
	StartDate   time.Time             `json:"start_date"`
	EndDate     time.Time             `json:"end_date"`
	Status      Status                `json:"status"`
	Progress    int                   `json:"progress"`
	Checkpoints []*CheckpointResponse `json:"checkpoints,omitempty"`
	Artifacts   []*ArtifactResponse   `json:"artifacts,omitempty"`
	Assessment  *AssessmentResponse   `json:"self_assessment,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type CheckpointResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     time.Time  `json:"due_date"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	OrderNum    int        `json:"order_num"`
}

type ArtifactResponse struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	Name string    `json:"name"`
	URL  string    `json:"url"`
}

type AssessmentResponse struct {
	ID      uuid.UUID `json:"id"`
	Rating  int       `json:"rating"`
	Comment string    `json:"comment"`
}

// ToResponse конвертирует Challenge в ChallengeResponse
func (c *Challenge) ToResponse(checkpoints []*Checkpoint, artifacts []*Artifact, assessment *SelfAssessment) *ChallengeResponse {
	resp := &ChallengeResponse{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		Goal:        c.Goal,
		StartDate:   c.StartDate,
		EndDate:     c.EndDate,
		Status:      c.Status,
		Progress:    c.Progress,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	// Добавляем чекпоинты
	if len(checkpoints) > 0 {
		resp.Checkpoints = make([]*CheckpointResponse, len(checkpoints))
		for i, cp := range checkpoints {
			resp.Checkpoints[i] = &CheckpointResponse{
				ID:          cp.ID,
				Title:       cp.Title,
				Description: cp.Description,
				DueDate:     cp.DueDate,
				IsCompleted: cp.IsCompleted,
				CompletedAt: cp.CompletedAt,
				OrderNum:    cp.OrderNum,
			}
		}
	}

	// Добавляем артефакты
	if len(artifacts) > 0 {
		resp.Artifacts = make([]*ArtifactResponse, len(artifacts))
		for i, a := range artifacts {
			resp.Artifacts[i] = &ArtifactResponse{
				ID:   a.ID,
				Type: a.Type,
				Name: a.Name,
				URL:  a.URL,
			}
		}
	}

	// Добавляем самооценку
	if assessment != nil {
		resp.Assessment = &AssessmentResponse{
			ID:      assessment.ID,
			Rating:  assessment.Rating,
			Comment: assessment.Comment,
		}
	}

	return resp
}

type CreateChallengeRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Goal        string    `json:"goal" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
}

type UpdateChallengeRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Goal        *string    `json:"goal,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Status      *Status    `json:"status,omitempty"`
	Progress    *int       `json:"progress,omitempty"`
}
