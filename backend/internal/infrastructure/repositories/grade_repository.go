package repositories

import (
	"context"
	"log"
	"time"

	"backend/internal/domain/grade"
	"backend/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GradeRepository struct {
	pool *pgxpool.Pool
}

func NewGradeRepository() *GradeRepository {
	return &GradeRepository{
		pool: db.Pool,
	}
}

// GetByUserID получает все оценки студента
func (r *GradeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*grade.Grade, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.user_id, g.course_id, g.type, g.value, g.max_value, 
		       g.weight, g.comment, g.date, g.source_type, g.source_id,
		       g.created_by, g.created_at, g.updated_at,
		       c.name as course_name
		FROM grades g
		JOIN courses c ON g.course_id = c.id
		WHERE g.user_id = $1
		ORDER BY g.date DESC
	`, userID)

	if err != nil {
		log.Printf("Database error in GetByUserID: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanGrades(rows)
}

// GetByUserAndCourse получает оценки по конкретному курсу
func (r *GradeRepository) GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) ([]*grade.Grade, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.user_id, g.course_id, g.type, g.value, g.max_value, 
		       g.weight, g.comment, g.date, g.source_type, g.source_id,
		       g.created_by, g.created_at, g.updated_at,
		       c.name as course_name
		FROM grades g
		JOIN courses c ON g.course_id = c.id
		WHERE g.user_id = $1 AND g.course_id = $2
		ORDER BY g.date DESC
	`, userID, courseID)

	if err != nil {
		log.Printf("Database error in GetByUserAndCourse: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanGrades(rows)
}

// GetByUserAndPeriod получает оценки за период
func (r *GradeRepository) GetByUserAndPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*grade.Grade, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.id, g.user_id, g.course_id, g.type, g.value, g.max_value, 
		       g.weight, g.comment, g.date, g.source_type, g.source_id,
		       g.created_by, g.created_at, g.updated_at,
		       c.name as course_name
		FROM grades g
		JOIN courses c ON g.course_id = c.id
		WHERE g.user_id = $1 AND g.date BETWEEN $2 AND $3
		ORDER BY g.date DESC
	`, userID, from, to)

	if err != nil {
		log.Printf("Database error in GetByUserAndPeriod: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanGrades(rows)
}

// GetByID получает оценку по ID
func (r *GradeRepository) GetByID(ctx context.Context, id uuid.UUID) (*grade.Grade, error) {
	g := &grade.Grade{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, course_id, type, value, max_value, 
		       weight, comment, date, source_type, source_id,
		       created_by, created_at, updated_at
		FROM grades
		WHERE id = $1
	`, id).Scan(
		&g.ID, &g.UserID, &g.CourseID, &g.Type, &g.Value, &g.MaxValue,
		&g.Weight, &g.Comment, &g.Date, &g.SourceType, &g.SourceID,
		&g.CreatedBy, &g.CreatedAt, &g.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByID: %v", err)
		return nil, err
	}
	return g, nil
}

// Create создаёт новую оценку
func (r *GradeRepository) Create(ctx context.Context, g *grade.Grade) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO grades (
			user_id, course_id, type, value, max_value, weight,
			comment, date, source_type, source_id, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`,
		g.UserID, g.CourseID, g.Type, g.Value, g.MaxValue, g.Weight,
		g.Comment, g.Date, g.SourceType, g.SourceID, g.CreatedBy,
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	if err != nil {
		log.Printf("Database error in Create: %v", err)
		return err
	}
	return nil
}

// CreateBatch создаёт несколько оценок (для импорта)
func (r *GradeRepository) CreateBatch(ctx context.Context, grades []*grade.Grade) error {
	batch := &pgx.Batch{}

	for _, g := range grades {
		batch.Queue(`
			INSERT INTO grades (
				user_id, course_id, type, value, max_value, weight,
				comment, date, source_type, source_id, created_by
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
			g.UserID, g.CourseID, g.Type, g.Value, g.MaxValue, g.Weight,
			g.Comment, g.Date, g.SourceType, g.SourceID, g.CreatedBy,
		)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	// Проверяем результаты
	for i := 0; i < len(grades); i++ {
		_, err := br.Exec()
		if err != nil {
			log.Printf("Error in batch insert at position %d: %v", i, err)
			return err
		}
	}

	return nil
}

// Update обновляет оценку
func (r *GradeRepository) Update(ctx context.Context, g *grade.Grade) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE grades
		SET value = $1, comment = $2, date = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, g.Value, g.Comment, g.Date, g.ID)

	if err != nil {
		log.Printf("Database error in Update: %v", err)
		return err
	}
	return nil
}

// Delete удаляет оценку
func (r *GradeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM grades WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Delete: %v", err)
		return err
	}
	return nil
}

// GetAverageByUser получает средний балл студента
func (r *GradeRepository) GetAverageByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	var avg float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(AVG(value / max_value * 100), 0)
		FROM grades
		WHERE user_id = $1
	`, userID).Scan(&avg)

	if err != nil {
		log.Printf("Database error in GetAverageByUser: %v", err)
		return 0, err
	}
	return avg, nil
}

// GetAverageByUserAndCourse получает средний балл по курсу
func (r *GradeRepository) GetAverageByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (float64, error) {
	var avg float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(AVG(value / max_value * 100), 0)
		FROM grades
		WHERE user_id = $1 AND course_id = $2
	`, userID, courseID).Scan(&avg)

	if err != nil {
		log.Printf("Database error in GetAverageByUserAndCourse: %v", err)
		return 0, err
	}
	return avg, nil
}

// GetSummaryByUser получает сводку по успеваемости студента
func (r *GradeRepository) GetSummaryByUser(ctx context.Context, userID uuid.UUID) (*grade.StudentSummary, error) {
	// Получаем общий средний балл
	overallAvg, err := r.GetAverageByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем сводку по курсам
	rows, err := r.pool.Query(ctx, `
		SELECT 
			g.course_id,
			c.name as course_name,
			COUNT(*) as grades_count,
			AVG(g.value / g.max_value * 100) as average,
			c.credits as total_credits
		FROM grades g
		JOIN courses c ON g.course_id = c.id
		WHERE g.user_id = $1
		GROUP BY g.course_id, c.name, c.credits
	`, userID)

	if err != nil {
		log.Printf("Database error in GetSummaryByUser: %v", err)
		return nil, err
	}
	defer rows.Close()

	var courses []*grade.CourseSummary
	totalCredits := 0

	for rows.Next() {
		cs := &grade.CourseSummary{}
		err := rows.Scan(
			&cs.CourseID, &cs.CourseName, &cs.GradesCount,
			&cs.Average, &cs.TotalCredits,
		)
		if err != nil {
			log.Printf("Error scanning course summary: %v", err)
			continue
		}
		courses = append(courses, cs)
		totalCredits += cs.TotalCredits
	}

	return &grade.StudentSummary{
		OverallAverage: overallAvg,
		Courses:        courses,
		TotalCredits:   totalCredits,
		LastUpdated:    time.Now(),
	}, nil
}

// Вспомогательная функция для сканирования строк
func (r *GradeRepository) scanGrades(rows pgx.Rows) ([]*grade.Grade, error) {
	var grades []*grade.Grade

	for rows.Next() {
		g := &grade.Grade{}
		var courseName string

		err := rows.Scan(
			&g.ID, &g.UserID, &g.CourseID, &g.Type, &g.Value, &g.MaxValue,
			&g.Weight, &g.Comment, &g.Date, &g.SourceType, &g.SourceID,
			&g.CreatedBy, &g.CreatedAt, &g.UpdatedAt,
			&courseName,
		)
		if err != nil {
			log.Printf("Error scanning grade row: %v", err)
			continue
		}
		grades = append(grades, g)
	}

	return grades, nil
}
