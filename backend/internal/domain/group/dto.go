package group

import (
	"github.com/google/uuid"
)

type CreateGroupRequest struct {
	Name       string     `json:"name" binding:"required"`
	CourseID   uuid.UUID  `json:"course_id" binding:"required"`
	SemesterID uuid.UUID  `json:"semester_id" binding:"required"`
	HolderID   *uuid.UUID `json:"holder_id,omitempty"`
}

type UpdateGroupRequest struct {
	Name       *string    `json:"name,omitempty"`
	CourseID   *uuid.UUID `json:"course_id,omitempty"`
	SemesterID *uuid.UUID `json:"semester_id,omitempty"`
	HolderID   *uuid.UUID `json:"holder_id,omitempty"`
}
