package application

import (
	"backend/internal/domain/group"
	"backend/internal/domain/semester"
	"context"
	"errors"

	"backend/internal/domain/user"

	"github.com/google/uuid"
)

type UserService struct {
	repo         user.Repository
	groupRepo    group.Repository
	semesterRepo semester.Repository
}

func NewUserService(repo user.Repository, groupRepo group.Repository, semesterRepo semester.Repository) *UserService {
	return &UserService{
		repo:         repo,
		groupRepo:    groupRepo,
		semesterRepo: semesterRepo,
	}
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

// GetProfileWithDetails получает профиль пользователя с информацией о группе и семестре
func (s *UserService) GetProfileWithDetails(ctx context.Context, id string) (*user.UserResponse, error) {
	// Получаем пользователя
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Создаём базовый ответ
	response := &user.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		GroupID:   u.GroupID,
	}

	// Если есть группа, получаем информацию о ней
	if u.GroupID != nil {
		grp, err := s.groupRepo.GetByID(ctx, *u.GroupID)
		if err == nil && grp != nil {
			response.GroupName = &grp.Name

			// Получаем информацию о семестре
			sem, err := s.semesterRepo.GetByID(ctx, grp.SemesterID)
			if err == nil && sem != nil {
				response.SemesterID = &sem.ID
				response.SemesterName = &sem.Name
				response.SemesterStart = &sem.StartDate
				response.SemesterEnd = &sem.EndDate
			}
		}
	}

	return response, nil
}
