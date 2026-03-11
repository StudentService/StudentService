package grade

import (
	"time"

	"github.com/google/uuid"
)

type GradeResponse struct {
	ID         uuid.UUID `json:"id"`
	CourseID   uuid.UUID `json:"course_id"`
	CourseName string    `json:"course_name"`
	Type       GradeType `json:"type"`
	Value      float64   `json:"value"`
	MaxValue   float64   `json:"max_value"`
	Weight     float64   `json:"weight"`
	Comment    string    `json:"comment"`
	Date       time.Time `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateGradeRequest struct {
	CourseID uuid.UUID `json:"course_id" binding:"required"`
	Type     GradeType `json:"type" binding:"required,oneof=exam test homework project activity"`
	Value    float64   `json:"value" binding:"required,min=0"`
	MaxValue float64   `json:"max_value" binding:"required,min=1"`
	Weight   float64   `json:"weight" binding:"min=0"`
	Comment  string    `json:"comment"`
	Date     time.Time `json:"date" binding:"required"`
}

type UpdateGradeRequest struct {
	Value   *float64   `json:"value,omitempty"`
	Comment *string    `json:"comment,omitempty"`
	Date    *time.Time `json:"date,omitempty"`
}

// ToResponse конвертирует Grade в GradeResponse
func (g *Grade) ToResponse(courseName string) *GradeResponse {
	return &GradeResponse{
		ID:         g.ID,
		CourseID:   g.CourseID,
		CourseName: courseName,
		Type:       g.Type,
		Value:      g.Value,
		MaxValue:   g.MaxValue,
		Weight:     g.Weight,
		Comment:    g.Comment,
		Date:       g.Date,
		CreatedAt:  g.CreatedAt,
	}
}
