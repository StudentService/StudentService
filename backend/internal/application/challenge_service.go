package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend/internal/domain/challenge"
	"backend/internal/domain/user"
)

type ChallengeService struct {
	challengeRepo challenge.Repository
	userRepo      user.Repository
}

func NewChallengeService(challengeRepo challenge.Repository, userRepo user.Repository) *ChallengeService {
	return &ChallengeService{
		challengeRepo: challengeRepo,
		userRepo:      userRepo,
	}
}

// GetMyChallenges получает все вызовы текущего пользователя
func (s *ChallengeService) GetMyChallenges(ctx context.Context, userID uuid.UUID) ([]*challenge.ChallengeResponse, error) {
	// Получаем вызовы пользователя
	challenges, err := s.challengeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Для каждого вызова получаем дополнительные данные
	responses := make([]*challenge.ChallengeResponse, len(challenges))
	for i, c := range challenges {
		// Обновляем статус просроченных вызовов
		s.updateChallengeStatus(c)

		// Получаем чекпоинты
		checkpoints, _ := s.challengeRepo.GetCheckpoints(ctx, c.ID)

		// Получаем артефакты
		artifacts, _ := s.challengeRepo.GetArtifacts(ctx, c.ID)

		// Получаем самооценку
		assessment, _ := s.challengeRepo.GetSelfAssessment(ctx, c.ID)

		responses[i] = c.ToResponse(checkpoints, artifacts, assessment)
	}

	return responses, nil
}

// GetChallengeByID получает конкретный вызов по ID
func (s *ChallengeService) GetChallengeByID(ctx context.Context, userID, challengeID uuid.UUID) (*challenge.ChallengeResponse, error) {
	// Получаем вызов
	c, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("challenge not found")
	}

	// Проверяем, что вызов принадлежит пользователю
	if c.UserID != userID {
		return nil, errors.New("access denied")
	}

	// Обновляем статус
	s.updateChallengeStatus(c)

	// Получаем связанные данные
	checkpoints, _ := s.challengeRepo.GetCheckpoints(ctx, c.ID)
	artifacts, _ := s.challengeRepo.GetArtifacts(ctx, c.ID)
	assessment, _ := s.challengeRepo.GetSelfAssessment(ctx, c.ID)

	return c.ToResponse(checkpoints, artifacts, assessment), nil
}

// CreateChallenge создает новый вызов
func (s *ChallengeService) CreateChallenge(ctx context.Context, userID uuid.UUID, req *challenge.CreateChallengeRequest) (*challenge.ChallengeResponse, error) {
	// Валидация дат
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Создаем вызов
	c := &challenge.Challenge{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Goal:        req.Goal,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      challenge.StatusActive,
		Progress:    0,
	}

	if err := s.challengeRepo.Create(ctx, c); err != nil {
		return nil, err
	}

	return c.ToResponse(nil, nil, nil), nil
}

// UpdateChallenge обновляет вызов
func (s *ChallengeService) UpdateChallenge(ctx context.Context, userID uuid.UUID, challengeID uuid.UUID, req *challenge.UpdateChallengeRequest) (*challenge.ChallengeResponse, error) {
	// Получаем существующий вызов
	c, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("challenge not found")
	}

	// Проверяем права
	if c.UserID != userID {
		return nil, errors.New("access denied")
	}

	// Обновляем поля
	if req.Title != nil {
		c.Title = *req.Title
	}
	if req.Description != nil {
		c.Description = *req.Description
	}
	if req.Goal != nil {
		c.Goal = *req.Goal
	}
	if req.StartDate != nil {
		c.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		c.EndDate = *req.EndDate
	}
	if req.Status != nil {
		c.Status = *req.Status
	}
	if req.Progress != nil {
		c.Progress = *req.Progress
	}

	// Сохраняем изменения
	if err := s.challengeRepo.Update(ctx, c); err != nil {
		return nil, err
	}

	// Получаем связанные данные для ответа
	checkpoints, _ := s.challengeRepo.GetCheckpoints(ctx, c.ID)
	artifacts, _ := s.challengeRepo.GetArtifacts(ctx, c.ID)
	assessment, _ := s.challengeRepo.GetSelfAssessment(ctx, c.ID)

	return c.ToResponse(checkpoints, artifacts, assessment), nil
}

// DeleteChallenge удаляет вызов
func (s *ChallengeService) DeleteChallenge(ctx context.Context, userID uuid.UUID, challengeID uuid.UUID) error {
	// Получаем вызов для проверки прав
	c, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("challenge not found")
	}

	// Проверяем права
	if c.UserID != userID {
		return errors.New("access denied")
	}

	return s.challengeRepo.Delete(ctx, challengeID)
}

// AddCheckpoint добавляет чекпоинт к вызову
func (s *ChallengeService) AddCheckpoint(ctx context.Context, userID uuid.UUID, challengeID uuid.UUID, req *challenge.Checkpoint) error {
	// Проверяем существование и права на вызов
	c, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}
	if c == nil || c.UserID != userID {
		return errors.New("access denied")
	}

	req.ChallengeID = challengeID
	return s.challengeRepo.CreateCheckpoint(ctx, req)
}

// UpdateCheckpoint обновляет чекпоинт
func (s *ChallengeService) UpdateCheckpoint(ctx context.Context, userID uuid.UUID, checkpointID uuid.UUID, isCompleted bool) error {
	// Получаем чекпоинт
	// Для этого нужен метод GetCheckpointByID, добавим позже при необходимости
	// Пока упростим - будем обновлять через challengeID
	return errors.New("not implemented")
}

// AddArtifact добавляет артефакт
func (s *ChallengeService) AddArtifact(ctx context.Context, userID uuid.UUID, challengeID uuid.UUID, req *challenge.Artifact) error {
	// Проверяем права
	c, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}
	if c == nil || c.UserID != userID {
		return errors.New("access denied")
	}

	req.ChallengeID = challengeID
	return s.challengeRepo.CreateArtifact(ctx, req)
}

// AddSelfAssessment добавляет самооценку
func (s *ChallengeService) AddSelfAssessment(ctx context.Context, userID uuid.UUID, challengeID uuid.UUID, rating int, comment string) error {
	// Проверяем права
	c, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}
	if c == nil || c.UserID != userID {
		return errors.New("access denied")
	}

	// Проверяем, нет ли уже самооценки
	existing, _ := s.challengeRepo.GetSelfAssessment(ctx, challengeID)
	if existing != nil {
		return errors.New("self assessment already exists")
	}

	assessment := &challenge.SelfAssessment{
		ChallengeID: challengeID,
		Rating:      rating,
		Comment:     comment,
	}

	return s.challengeRepo.CreateSelfAssessment(ctx, assessment)
}

// updateChallengeStatus обновляет статус вызова если он просрочен
func (s *ChallengeService) updateChallengeStatus(c *challenge.Challenge) {
	if c.Status != challenge.StatusCompleted && c.Status != challenge.StatusOverdue {
		if time.Now().After(c.EndDate) {
			c.Status = challenge.StatusOverdue
			// Здесь можно вызвать repo.Update, но это отдельный запрос
			// Для простоты пока не обновляем в БД автоматически
		}
	}
}
