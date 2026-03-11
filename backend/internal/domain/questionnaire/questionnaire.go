package questionnaire

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSON тип для Swagger
type JSON map[string]interface{}

// Или более специфичный тип
type AnswersJSON map[string]interface{}

type Status string

const (
	StatusDraft     Status = "draft"     // черновик
	StatusSubmitted Status = "submitted" // отправлено
	StatusApproved  Status = "approved"  // одобрено
	StatusRejected  Status = "rejected"  // отклонено
)

// Questionnaire - анкета студента
type Questionnaire struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	Status      Status          `json:"status"`
	Answers     json.RawMessage `json:"answers" swaggertype:"object"` // добавляем swaggertype
	SubmittedAt *time.Time      `json:"submitted_at,omitempty"`
	ReviewedBy  *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt  *time.Time      `json:"reviewed_at,omitempty"`
	Comment     string          `json:"comment,omitempty"` // комментарий при отклонении
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// QuestionnaireTemplate - шаблон анкеты (настраивается админом)
type QuestionnaireTemplate struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	Schema      json.RawMessage `json:"schema" swaggertype:"object"`
	Fields      json.RawMessage `json:"fields" swaggertype:"object"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// FieldType типы полей анкеты
type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeTextarea FieldType = "textarea"
	FieldTypeNumber   FieldType = "number"
	FieldTypeSelect   FieldType = "select"
	FieldTypeRadio    FieldType = "radio"
	FieldTypeCheckbox FieldType = "checkbox"
	FieldTypeDate     FieldType = "date"
	FieldTypeFile     FieldType = "file"
)

// FieldDefinition определение поля анкеты
type FieldDefinition struct {
	ID          string                 `json:"id"`
	Type        FieldType              `json:"type"`
	Label       string                 `json:"label"`
	Required    bool                   `json:"required"`
	Placeholder string                 `json:"placeholder,omitempty"`
	Options     []FieldOption          `json:"options,omitempty"`    // для select/radio
	Validation  map[string]interface{} `json:"validation,omitempty"` // правила валидации
}

// FieldOption вариант для выбора
type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Repository interface {
	// Для студента
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Questionnaire, error)
	Create(ctx context.Context, q *Questionnaire) error
	Update(ctx context.Context, q *Questionnaire) error

	// Для админа (управление шаблонами)
	GetActiveTemplate(ctx context.Context) (*QuestionnaireTemplate, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*QuestionnaireTemplate, error)
	CreateTemplate(ctx context.Context, t *QuestionnaireTemplate) error
	UpdateTemplate(ctx context.Context, t *QuestionnaireTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	ListTemplates(ctx context.Context) ([]*QuestionnaireTemplate, error)

	// Для админа (просмотр анкет)
	ListByStatus(ctx context.Context, status Status) ([]*Questionnaire, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Questionnaire, error)
	Review(ctx context.Context, id uuid.UUID, status Status, reviewerID uuid.UUID, comment string) error
}
