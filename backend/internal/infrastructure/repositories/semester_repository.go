package repositories

import (
	"context"
	"log"

	"backend/internal/domain/semester"
	"backend/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SemesterRepository struct{}

func (r *SemesterRepository) GetByID(ctx context.Context, id uuid.UUID) (*semester.Semester, error) {
	s := &semester.Semester{}
	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, start_date, end_date, is_active, created_at, updated_at 
		FROM semesters WHERE id = $1
	`, id).Scan(
		&s.ID, &s.Name, &s.StartDate, &s.EndDate, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in Semester.GetByID: %v", err)
		return nil, err
	}
	return s, nil
}

func (r *SemesterRepository) GetActive(ctx context.Context) (*semester.Semester, error) {
	s := &semester.Semester{}
	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, start_date, end_date, is_active, created_at, updated_at 
		FROM semesters WHERE is_active = true LIMIT 1
	`).Scan(
		&s.ID, &s.Name, &s.StartDate, &s.EndDate, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in Semester.GetActive: %v", err)
		return nil, err
	}
	return s, nil
}

func (r *SemesterRepository) GetAll(ctx context.Context) ([]*semester.Semester, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, start_date, end_date, is_active, created_at, updated_at 
		FROM semesters ORDER BY start_date DESC
	`)
	if err != nil {
		log.Printf("Database error in Semester.GetAll: %v", err)
		return nil, err
	}
	defer rows.Close()

	var semesters []*semester.Semester
	for rows.Next() {
		s := &semester.Semester{}
		err := rows.Scan(
			&s.ID, &s.Name, &s.StartDate, &s.EndDate, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning semester row: %v", err)
			continue
		}
		semesters = append(semesters, s)
	}

	return semesters, nil
}

func (r *SemesterRepository) Create(ctx context.Context, s *semester.Semester) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO semesters (name, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, s.Name, s.StartDate, s.EndDate, s.IsActive).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		log.Printf("Database error in Semester.Create: %v", err)
		return err
	}
	return nil
}

func (r *SemesterRepository) Update(ctx context.Context, s *semester.Semester) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE semesters
		SET name = $1, start_date = $2, end_date = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`, s.Name, s.StartDate, s.EndDate, s.IsActive, s.ID)

	if err != nil {
		log.Printf("Database error in Semester.Update: %v", err)
		return err
	}
	return nil
}

func (r *SemesterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM semesters WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Semester.Delete: %v", err)
		return err
	}
	return nil
}
