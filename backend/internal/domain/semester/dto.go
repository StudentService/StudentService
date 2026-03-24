package semester

import (
	"time"
)

type CreateSemesterRequest struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	IsActive  bool      `json:"is_active"`
}

type UpdateSemesterRequest struct {
	Name      *string    `json:"name,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}
