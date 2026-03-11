package repositories

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"backend/internal/domain/questionnaire"
	"backend/internal/infrastructure/db"
)

type QuestionnaireRepository struct{}

// GetByUserID получает анкету студента по его ID
func (r *QuestionnaireRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*questionnaire.Questionnaire, error) {
	q := &questionnaire.Questionnaire{}
	var reviewedBy *uuid.UUID
	var answers []byte

	err := db.Pool.QueryRow(ctx, `
		SELECT id, user_id, status, answers, submitted_at, 
		       reviewed_by, reviewed_at, comment, created_at, updated_at
		FROM questionnaires
		WHERE user_id = $1
	`, userID).Scan(
		&q.ID, &q.UserID, &q.Status, &answers, &q.SubmittedAt,
		&reviewedBy, &q.ReviewedAt, &q.Comment, &q.CreatedAt, &q.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByUserID: %v", err)
		return nil, err
	}

	q.Answers = answers
	q.ReviewedBy = reviewedBy
	return q, nil
}

// Create создаёт новую анкету
func (r *QuestionnaireRepository) Create(ctx context.Context, q *questionnaire.Questionnaire) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO questionnaires (user_id, status, answers)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, q.UserID, q.Status, q.Answers).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)

	if err != nil {
		log.Printf("Database error in Create: %v", err)
		return err
	}
	return nil
}

// Update обновляет анкету
func (r *QuestionnaireRepository) Update(ctx context.Context, q *questionnaire.Questionnaire) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE questionnaires
		SET status = $1, answers = $2, submitted_at = $3,
		    reviewed_by = $4, reviewed_at = $5, comment = $6,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`, q.Status, q.Answers, q.SubmittedAt, q.ReviewedBy, q.ReviewedAt, q.Comment, q.ID)

	if err != nil {
		log.Printf("Database error in Update: %v", err)
		return err
	}
	return nil
}

// GetActiveTemplate получает активный шаблон анкеты
func (r *QuestionnaireRepository) GetActiveTemplate(ctx context.Context) (*questionnaire.QuestionnaireTemplate, error) {
	t := &questionnaire.QuestionnaireTemplate{}
	var schema, fields []byte

	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, description, is_active, schema, fields, created_at, updated_at
		FROM questionnaire_templates
		WHERE is_active = true
		LIMIT 1
	`).Scan(
		&t.ID, &t.Name, &t.Description, &t.IsActive, &schema, &fields,
		&t.CreatedAt, &t.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetActiveTemplate: %v", err)
		return nil, err
	}

	t.Schema = schema
	t.Fields = fields
	return t, nil
}

// GetTemplateByID получает шаблон по ID
func (r *QuestionnaireRepository) GetTemplateByID(ctx context.Context, id uuid.UUID) (*questionnaire.QuestionnaireTemplate, error) {
	t := &questionnaire.QuestionnaireTemplate{}
	var schema, fields []byte

	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, description, is_active, schema, fields, created_at, updated_at
		FROM questionnaire_templates
		WHERE id = $1
	`, id).Scan(
		&t.ID, &t.Name, &t.Description, &t.IsActive, &schema, &fields,
		&t.CreatedAt, &t.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetTemplateByID: %v", err)
		return nil, err
	}

	t.Schema = schema
	t.Fields = fields
	return t, nil
}

// CreateTemplate создаёт новый шаблон
func (r *QuestionnaireRepository) CreateTemplate(ctx context.Context, t *questionnaire.QuestionnaireTemplate) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO questionnaire_templates (name, description, is_active, schema, fields)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, t.Name, t.Description, t.IsActive, t.Schema, t.Fields).Scan(
		&t.ID, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		log.Printf("Database error in CreateTemplate: %v", err)
		return err
	}
	return nil
}

// UpdateTemplate обновляет шаблон
func (r *QuestionnaireRepository) UpdateTemplate(ctx context.Context, t *questionnaire.QuestionnaireTemplate) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE questionnaire_templates
		SET name = $1, description = $2, is_active = $3, 
		    schema = $4, fields = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`, t.Name, t.Description, t.IsActive, t.Schema, t.Fields, t.ID)

	if err != nil {
		log.Printf("Database error in UpdateTemplate: %v", err)
		return err
	}
	return nil
}

// DeleteTemplate удаляет шаблон
func (r *QuestionnaireRepository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM questionnaire_templates WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in DeleteTemplate: %v", err)
		return err
	}
	return nil
}

// ListTemplates получает все шаблоны
func (r *QuestionnaireRepository) ListTemplates(ctx context.Context) ([]*questionnaire.QuestionnaireTemplate, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, description, is_active, schema, fields, created_at, updated_at
		FROM questionnaire_templates
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Database error in ListTemplates: %v", err)
		return nil, err
	}
	defer rows.Close()

	var templates []*questionnaire.QuestionnaireTemplate
	for rows.Next() {
		t := &questionnaire.QuestionnaireTemplate{}
		var schema, fields []byte

		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.IsActive, &schema, &fields,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning template row: %v", err)
			continue
		}
		t.Schema = schema
		t.Fields = fields
		templates = append(templates, t)
	}

	return templates, nil
}

// ListByStatus получает анкеты по статусу (для админа)
func (r *QuestionnaireRepository) ListByStatus(ctx context.Context, status questionnaire.Status) ([]*questionnaire.Questionnaire, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, user_id, status, answers, submitted_at, 
		       reviewed_by, reviewed_at, comment, created_at, updated_at
		FROM questionnaires
		WHERE status = $1
		ORDER BY submitted_at DESC NULLS LAST, created_at DESC
	`, status)

	if err != nil {
		log.Printf("Database error in ListByStatus: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanQuestionnaires(rows)
}

// GetByID получает анкету по ID (для админа)
func (r *QuestionnaireRepository) GetByID(ctx context.Context, id uuid.UUID) (*questionnaire.Questionnaire, error) {
	q := &questionnaire.Questionnaire{}
	var reviewedBy *uuid.UUID
	var answers []byte

	err := db.Pool.QueryRow(ctx, `
		SELECT id, user_id, status, answers, submitted_at, 
		       reviewed_by, reviewed_at, comment, created_at, updated_at
		FROM questionnaires
		WHERE id = $1
	`, id).Scan(
		&q.ID, &q.UserID, &q.Status, &answers, &q.SubmittedAt,
		&reviewedBy, &q.ReviewedAt, &q.Comment, &q.CreatedAt, &q.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByID: %v", err)
		return nil, err
	}

	q.Answers = answers
	q.ReviewedBy = reviewedBy
	return q, nil
}

// Review обновляет статус проверки анкеты
func (r *QuestionnaireRepository) Review(ctx context.Context, id uuid.UUID, status questionnaire.Status, reviewerID uuid.UUID, comment string) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE questionnaires
		SET status = $1, reviewed_by = $2, reviewed_at = CURRENT_TIMESTAMP, 
		    comment = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, status, reviewerID, comment, id)

	if err != nil {
		log.Printf("Database error in Review: %v", err)
		return err
	}
	return nil
}

// Вспомогательная функция для сканирования анкет
func (r *QuestionnaireRepository) scanQuestionnaires(rows pgx.Rows) ([]*questionnaire.Questionnaire, error) {
	var questionnaires []*questionnaire.Questionnaire

	for rows.Next() {
		q := &questionnaire.Questionnaire{}
		var reviewedBy *uuid.UUID
		var answers []byte

		err := rows.Scan(
			&q.ID, &q.UserID, &q.Status, &answers, &q.SubmittedAt,
			&reviewedBy, &q.ReviewedAt, &q.Comment, &q.CreatedAt, &q.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning questionnaire row: %v", err)
			continue
		}

		q.Answers = answers
		q.ReviewedBy = reviewedBy
		questionnaires = append(questionnaires, q)
	}

	return questionnaires, nil
}
