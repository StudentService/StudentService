package repositories

import (
	"context"
	"log"

	"backend/internal/domain/course"
	"backend/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CourseRepository struct{}

func (r *CourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*course.Course, error) {
	c := &course.Course{}
	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, description, credits, created_at, updated_at 
		FROM courses WHERE id = $1
	`, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.Credits, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in Course.GetByID: %v", err)
		return nil, err
	}
	return c, nil
}

func (r *CourseRepository) GetAll(ctx context.Context) ([]*course.Course, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, description, credits, created_at, updated_at 
		FROM courses ORDER BY name
	`)
	if err != nil {
		log.Printf("Database error in Course.GetAll: %v", err)
		return nil, err
	}
	defer rows.Close()

	var courses []*course.Course
	for rows.Next() {
		c := &course.Course{}
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.Credits, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning course row: %v", err)
			continue
		}
		courses = append(courses, c)
	}

	return courses, nil
}

func (r *CourseRepository) Create(ctx context.Context, c *course.Course) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO courses (name, description, credits)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, c.Name, c.Description, c.Credits).Scan(
		&c.ID, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		log.Printf("Database error in Course.Create: %v", err)
		return err
	}
	return nil
}

func (r *CourseRepository) Update(ctx context.Context, c *course.Course) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE courses
		SET name = $1, description = $2, credits = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, c.Name, c.Description, c.Credits, c.ID)

	if err != nil {
		log.Printf("Database error in Course.Update: %v", err)
		return err
	}
	return nil
}

func (r *CourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM courses WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Course.Delete: %v", err)
		return err
	}
	return nil
}
