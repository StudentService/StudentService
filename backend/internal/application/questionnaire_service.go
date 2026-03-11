package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"

	"backend/internal/domain/questionnaire"
	"backend/internal/domain/user"
)

type QuestionnaireService struct {
	questionnaireRepo questionnaire.Repository
	userRepo          user.Repository
}

func NewQuestionnaireService(
	questionnaireRepo questionnaire.Repository,
	userRepo user.Repository,
) *QuestionnaireService {
	return &QuestionnaireService{
		questionnaireRepo: questionnaireRepo,
		userRepo:          userRepo,
	}
}

// GetMyQuestionnaire получает анкету текущего пользователя
func (s *QuestionnaireService) GetMyQuestionnaire(ctx context.Context, userID uuid.UUID) (*questionnaire.QuestionnaireResponse, error) {
	// Получаем анкету
	q, err := s.questionnaireRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Если анкеты нет, возвращаем пустой ответ с инструкцией создать
	if q == nil {
		return nil, nil
	}

	return q.ToResponse(), nil
}

// GetTemplate получает активный шаблон анкеты
func (s *QuestionnaireService) GetTemplate(ctx context.Context) (*questionnaire.TemplateResponse, error) {
	template, err := s.questionnaireRepo.GetActiveTemplate(ctx)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errors.New("no active questionnaire template found")
	}

	return template.ToTemplateResponse(), nil
}

// SubmitQuestionnaire отправляет заполненную анкету
func (s *QuestionnaireService) SubmitQuestionnaire(ctx context.Context, userID uuid.UUID, req *questionnaire.SubmitRequest) (*questionnaire.SubmitResponse, error) {
	// Получаем активный шаблон для валидации
	template, err := s.questionnaireRepo.GetActiveTemplate(ctx)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errors.New("no active questionnaire template found")
	}

	// Конвертируем map в json.RawMessage для валидации
	answersJSON, err := json.Marshal(req.Answers)
	if err != nil {
		return nil, errors.New("failed to marshal answers")
	}

	// Валидируем ответы по схеме
	if err := s.validateAnswers(template.Schema, answersJSON); err != nil {
		return nil, err
	}

	// Проверяем, есть ли уже анкета у пользователя
	existing, err := s.questionnaireRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if existing == nil {
		// Создаём новую анкету
		q := &questionnaire.Questionnaire{
			UserID:      userID,
			Status:      questionnaire.StatusSubmitted,
			Answers:     answersJSON,
			SubmittedAt: &now,
		}

		if err := s.questionnaireRepo.Create(ctx, q); err != nil {
			return nil, err
		}

		return &questionnaire.SubmitResponse{
			ID:          q.ID,
			Status:      q.Status,
			SubmittedAt: q.SubmittedAt,
		}, nil
	} else {
		// Обновляем существующую
		if existing.Status != questionnaire.StatusDraft && existing.Status != questionnaire.StatusRejected {
			return nil, errors.New("cannot update questionnaire in current status")
		}

		existing.Status = questionnaire.StatusSubmitted
		existing.Answers = answersJSON
		existing.SubmittedAt = &now
		existing.ReviewedBy = nil
		existing.ReviewedAt = nil
		existing.Comment = ""

		if err := s.questionnaireRepo.Update(ctx, existing); err != nil {
			return nil, err
		}

		return &questionnaire.SubmitResponse{
			ID:          existing.ID,
			Status:      existing.Status,
			SubmittedAt: existing.SubmittedAt,
		}, nil
	}
}

// SaveDraft сохраняет черновик анкеты
func (s *QuestionnaireService) SaveDraft(ctx context.Context, userID uuid.UUID, answers map[string]interface{}) error {
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return errors.New("failed to marshal answers")
	}

	// Проверяем, есть ли уже анкета
	existing, err := s.questionnaireRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Создаём черновик
		q := &questionnaire.Questionnaire{
			UserID:  userID,
			Status:  questionnaire.StatusDraft,
			Answers: answersJSON,
		}
		return s.questionnaireRepo.Create(ctx, q)
	} else {
		// Обновляем черновик
		if existing.Status != questionnaire.StatusDraft && existing.Status != questionnaire.StatusRejected {
			return errors.New("cannot edit questionnaire in current status")
		}

		existing.Answers = answersJSON
		existing.UpdatedAt = time.Now()
		return s.questionnaireRepo.Update(ctx, existing)
	}
}

// ReviewQuestionnaire проверяет анкету (для админа)
func (s *QuestionnaireService) ReviewQuestionnaire(ctx context.Context, reviewerID uuid.UUID, questionnaireID uuid.UUID, req *questionnaire.ReviewRequest) error {
	// Проверяем права рецензента
	reviewer, err := s.userRepo.GetByID(ctx, reviewerID.String())
	if err != nil {
		return err
	}
	if reviewer == nil {
		return errors.New("reviewer not found")
	}

	// Только админ может проверять анкеты
	if !reviewer.IsAdmin() {
		return errors.New("insufficient permissions")
	}

	// Получаем анкету
	q, err := s.questionnaireRepo.GetByID(ctx, questionnaireID)
	if err != nil {
		return err
	}
	if q == nil {
		return errors.New("questionnaire not found")
	}

	// Можно проверять только отправленные анкеты
	if q.Status != questionnaire.StatusSubmitted {
		return errors.New("can only review submitted questionnaires")
	}

	// Обновляем статус
	return s.questionnaireRepo.Review(ctx, questionnaireID, req.Status, reviewerID, req.Comment)
}

// ListQuestionnairesByStatus получает анкеты по статусу (для админа)
func (s *QuestionnaireService) ListQuestionnairesByStatus(ctx context.Context, status questionnaire.Status) ([]*questionnaire.QuestionnaireResponse, error) {
	questionnaires, err := s.questionnaireRepo.ListByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	responses := make([]*questionnaire.QuestionnaireResponse, len(questionnaires))
	for i, q := range questionnaires {
		responses[i] = q.ToResponse()
	}

	return responses, nil
}

// validateAnswers проверяет ответы по JSON схеме
func (s *QuestionnaireService) validateAnswers(schemaJSON json.RawMessage, answersJSON json.RawMessage) error {
	// Загружаем схему
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return errors.New("invalid schema: " + err.Error())
	}

	// Загружаем ответы
	answersLoader := gojsonschema.NewBytesLoader(answersJSON)

	// Валидируем
	result, err := schema.Validate(answersLoader)
	if err != nil {
		return errors.New("validation error: " + err.Error())
	}

	if !result.Valid() {
		// Собираем ошибки
		errs := ""
		for i, desc := range result.Errors() {
			if i > 0 {
				errs += "; "
			}
			errs += desc.String()
		}
		return errors.New("validation failed: " + errs)
	}

	return nil
}
