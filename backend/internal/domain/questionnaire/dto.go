package questionnaire

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SubmitRequest DTO для отправки анкеты
type SubmitRequest struct {
	Answers map[string]interface{} `json:"answers" binding:"required"` // меняем на map
}

// SubmitResponse DTO для ответа после отправки
type SubmitResponse struct {
	ID          uuid.UUID  `json:"id"`
	Status      Status     `json:"status"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}

// QuestionnaireResponse DTO для получения анкеты
type QuestionnaireResponse struct {
	ID          uuid.UUID              `json:"id"`
	Status      Status                 `json:"status"`
	Answers     map[string]interface{} `json:"answers"`
	SubmittedAt *time.Time             `json:"submitted_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`

	// Поля для проверки (если есть)
	ReviewedBy *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	Comment    string     `json:"comment,omitempty"`
}

// TemplateResponse DTO для шаблона анкеты - ИСПРАВЛЕНО
type TemplateResponse struct {
	ID          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Fields      []map[string]interface{} `json:"fields"` // меняем на slice of maps
	IsActive    bool                     `json:"is_active"`
}

// ReviewRequest DTO для проверки анкеты (админ)
type ReviewRequest struct {
	Status  Status `json:"status" binding:"required,oneof=approved rejected"`
	Comment string `json:"comment,omitempty"`
}

// ToResponse конвертирует Questionnaire в QuestionnaireResponse
func (q *Questionnaire) ToResponse() *QuestionnaireResponse {
	var answersMap map[string]interface{}
	if q.Answers != nil {
		json.Unmarshal(q.Answers, &answersMap)
	}

	return &QuestionnaireResponse{
		ID:          q.ID,
		Status:      q.Status,
		Answers:     answersMap,
		SubmittedAt: q.SubmittedAt,
		CreatedAt:   q.CreatedAt,
		UpdatedAt:   q.UpdatedAt,
		ReviewedBy:  q.ReviewedBy,
		ReviewedAt:  q.ReviewedAt,
		Comment:     q.Comment,
	}
}

// ToTemplateResponse конвертирует QuestionnaireTemplate в TemplateResponse - ИСПРАВЛЕНО
func (t *QuestionnaireTemplate) ToTemplateResponse() *TemplateResponse {
	var fields []map[string]interface{}
	if t.Fields != nil {
		json.Unmarshal(t.Fields, &fields)
	}

	return &TemplateResponse{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Fields:      fields,
		IsActive:    t.IsActive,
	}
}
