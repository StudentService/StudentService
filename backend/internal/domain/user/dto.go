package user

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse DTO для безопасной отправки данных пользователя
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      Role       `json:"role"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	GroupID   *uuid.UUID `json:"group_id,omitempty"`
	HolderID  *uuid.UUID `json:"holder_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
