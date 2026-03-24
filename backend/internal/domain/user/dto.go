package user

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse DTO для безопасной отправки данных пользователя
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Информация о группе
	GroupID   *uuid.UUID `json:"group_id,omitempty"`
	GroupName *string    `json:"group_name,omitempty"`

	// Информация о семестре
	SemesterID    *uuid.UUID `json:"semester_id,omitempty"`
	SemesterName  *string    `json:"semester_name,omitempty"`
	SemesterStart *time.Time `json:"semester_start,omitempty"`
	SemesterEnd   *time.Time `json:"semester_end,omitempty"`
}

// ToResponse конвертирует User в UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		GroupID:   u.GroupID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
