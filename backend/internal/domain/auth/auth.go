package auth

import (
	"backend/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // в секундах
}

// RegisterRequest DTO для регистрации
type RegisterRequest struct {
	Username  string     `json:"username" binding:"required,min=3,max=50"`
	Email     string     `json:"email" binding:"required,email"`
	Password  string     `json:"password" binding:"required,min=6"`
	FirstName string     `json:"first_name" binding:"required"`
	LastName  string     `json:"last_name" binding:"required"`
	Role      string     `json:"role" binding:"omitempty,oneof=student teacher holder candidate"`
	GroupID   *uuid.UUID `json:"group_id,omitempty"`
}

type RegisterResponse struct {
	User  *user.UserResponse `json:"user"`
	Token *AuthResponse      `json:"token,omitempty"`
}
