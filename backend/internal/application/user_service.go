package application

import (
	"context"
	"errors"

	"backend/internal/domain/user"

	"github.com/google/uuid"
)

type UserService struct {
	repo user.Repository
}

func NewUserService(repo user.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(ctx context.Context, id string) (*user.User, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, u *user.User) error {
	if u == nil || u.ID == uuid.Nil {
		return errors.New("invalid user data")
	}

	return s.repo.Update(ctx, u)
}

func (s *UserService) CreateUser(ctx context.Context, userData *user.User, password string) (*user.User, error) {
	// БИЗНЕС-ЛОГИКА (это слой application)

	// 1. Устанавливаем пароль через метод модели
	if err := userData.SetPassword(password); err != nil {
		return nil, err
	}

	// 2. Валидация данных
	if userData.Email == "" {
		return nil, errors.New("email is required")
	}

	// 3. Проверка уникальности (бизнес-правило)
	existingUser, _ := s.repo.GetByEmail(ctx, userData.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 4. Сохраняем в БД через ВАШ РЕПОЗИТОРИЙ
	return s.repo.Create(ctx, userData) // ← вызываем ваш метод
}
