package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleStudent   Role = "student"
	RoleTeacher   Role = "teacher"
	RoleHolder    Role = "holder"
	RoleCandidate Role = "candidate"
	RoleAdmin     Role = "admin"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SetPassword хеширует и устанавливает пароль
func (u *User) SetPassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword проверяет соответствие пароля
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Методы для проверки ролей
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

func (u *User) IsHolder() bool {
	return u.Role == RoleHolder
}

func (u *User) IsCandidate() bool {
	return u.Role == RoleCandidate
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// Методы для проверки прав доступа (временная упрощенная версия)
func (u *User) CanViewUser(target *User) bool {
	switch u.Role {
	case RoleAdmin:
		// Админ видит всех
		return true
	case RoleHolder:
		// Держатель пока видит всех (временное решение)
		// TODO: Implement holder-student relationship when HolderID is added
		return true
	case RoleTeacher:
		// Преподаватель пока видит всех (временное решение)
		// TODO: Implement teacher-group relationship
		return true
	case RoleStudent:
		// Студент видит только себя
		return u.ID == target.ID
	default:
		return false
	}
}

func (u *User) CanEditUser(target *User) bool {
	switch u.Role {
	case RoleAdmin:
		// Админ может редактировать всех
		return true
	case RoleHolder:
		// Держатель пока может редактировать всех (временное решение)
		// TODO: Restrict to only their students when HolderID is added
		return true
	case RoleTeacher:
		// Преподаватель пока не может редактировать профили
		// (только оценки, но это будет в других методах)
		return false
	case RoleStudent:
		// Студент может редактировать только свой профиль
		return u.ID == target.ID
	default:
		return false
	}
}

// Вспомогательные методы
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsActive() bool {
	// Здесь можно добавить логику проверки активности
	// Например, не забанен, подтвержден email и т.д.
	return true
}
