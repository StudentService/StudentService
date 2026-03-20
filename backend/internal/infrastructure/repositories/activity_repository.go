package repositories

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/activity"
	"backend/internal/infrastructure/db"
)

type ActivityRepository struct {
	pool *pgxpool.Pool
}

func NewActivityRepository() *ActivityRepository {
	return &ActivityRepository{
		pool: db.Pool,
	}
}

// GetAvailableActivities получает активности, доступные для студента
func (r *ActivityRepository) GetAvailableActivities(ctx context.Context, userID uuid.UUID) ([]*activity.Activity, error) {
	query := `
		SELECT DISTINCT a.id, a.title, a.description, a.type, a.status,
		       a.start_time, a.end_time, a.deadline,
		       a.location, a.online_link,
		       a.max_participants, a.current_participants,
		       a.points, a.weight,
		       a.course_id, a.group_id,
		       a.created_by, a.created_by_role,
		       a.created_at, a.updated_at,
		       CASE WHEN p.id IS NOT NULL THEN true ELSE false END as is_enrolled,
		       p.status as enrollment_status
		FROM activities a
		LEFT JOIN participations p ON a.id = p.activity_id AND p.user_id = $1
		LEFT JOIN users u ON u.id = $1
		WHERE a.status = 'active'
		  AND (
		      a.course_id IS NULL OR
		      a.group_id IS NULL OR
		      a.group_id = u.group_id OR
		      a.course_id IN (
		          SELECT course_id FROM groups WHERE id = u.group_id
		      )
		  )
		ORDER BY a.start_time NULLS LAST, a.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Database error in GetAvailableActivities: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

// GetByID получает активность по ID
func (r *ActivityRepository) GetByID(ctx context.Context, id uuid.UUID) (*activity.Activity, error) {
	a := &activity.Activity{}
	var courseID, groupID *uuid.UUID

	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, type, status,
		       start_time, end_time, deadline,
		       location, online_link,
		       max_participants, current_participants,
		       points, weight,
		       course_id, group_id,
		       created_by, created_by_role,
		       created_at, updated_at
		FROM activities
		WHERE id = $1
	`, id).Scan(
		&a.ID, &a.Title, &a.Description, &a.Type, &a.Status,
		&a.StartTime, &a.EndTime, &a.Deadline,
		&a.Location, &a.OnlineLink,
		&a.MaxParticipants, &a.CurrentParticipants,
		&a.Points, &a.Weight,
		&courseID, &groupID,
		&a.CreatedBy, &a.CreatedByRole,
		&a.CreatedAt, &a.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByID: %v", err)
		return nil, err
	}

	a.CourseID = courseID
	a.GroupID = groupID
	return a, nil
}

// GetMyParticipations получает все участия студента
func (r *ActivityRepository) GetMyParticipations(ctx context.Context, userID uuid.UUID) ([]*activity.Participation, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.activity_id, p.user_id, p.status,
		       p.grade, p.feedback, p.points_earned,
		       p.enrolled_at, p.completed_at,
		       p.created_at, p.updated_at,
		       a.title, a.type, a.start_time, a.location
		FROM participations p
		JOIN activities a ON p.activity_id = a.id
		WHERE p.user_id = $1
		ORDER BY p.enrolled_at DESC
	`, userID)

	if err != nil {
		log.Printf("Database error in GetMyParticipations: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanParticipations(rows)
}

// GetParticipationByActivity получает участие студента в конкретной активности
func (r *ActivityRepository) GetParticipationByActivity(ctx context.Context, userID, activityID uuid.UUID) (*activity.Participation, error) {
	p := &activity.Participation{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, activity_id, user_id, status,
		       grade, feedback, points_earned,
		       enrolled_at, completed_at,
		       created_at, updated_at
		FROM participations
		WHERE user_id = $1 AND activity_id = $2
	`, userID, activityID).Scan(
		&p.ID, &p.ActivityID, &p.UserID, &p.Status,
		&p.Grade, &p.Feedback, &p.PointsEarned,
		&p.EnrolledAt, &p.CompletedAt,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetParticipationByActivity: %v", err)
		return nil, err
	}

	return p, nil
}

// Enroll записывает студента на активность
func (r *ActivityRepository) Enroll(ctx context.Context, p *activity.Participation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Проверяем, есть ли места
	var currentParticipants, maxParticipants int
	err = tx.QueryRow(ctx, `
		SELECT current_participants, max_participants
		FROM activities
		WHERE id = $1
	`, p.ActivityID).Scan(&currentParticipants, &maxParticipants)
	if err != nil {
		return err
	}

	if maxParticipants > 0 && currentParticipants >= maxParticipants {
		return pgx.ErrNoRows // специально используем для сигнала "нет мест"
	}

	// Создаём участие
	err = tx.QueryRow(ctx, `
		INSERT INTO participations (
			activity_id, user_id, status, points_earned
		) VALUES ($1, $2, $3, $4)
		RETURNING id, enrolled_at, created_at, updated_at
	`, p.ActivityID, p.UserID, p.Status, p.PointsEarned).Scan(
		&p.ID, &p.EnrolledAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Увеличиваем счётчик участников
	_, err = tx.Exec(ctx, `
		UPDATE activities
		SET current_participants = current_participants + 1
		WHERE id = $1
	`, p.ActivityID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// CancelEnrollment отменяет запись на активность
func (r *ActivityRepository) CancelEnrollment(ctx context.Context, id uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Получаем activity_id перед удалением
	var activityID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT activity_id FROM participations WHERE id = $1
	`, id).Scan(&activityID)
	if err != nil {
		return err
	}

	// Удаляем участие
	_, err = tx.Exec(ctx, `DELETE FROM participations WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// Уменьшаем счётчик участников
	_, err = tx.Exec(ctx, `
		UPDATE activities
		SET current_participants = current_participants - 1
		WHERE id = $1 AND current_participants > 0
	`, activityID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// UpdateParticipationStatus обновляет статус участия
func (r *ActivityRepository) UpdateParticipationStatus(ctx context.Context, id uuid.UUID, status activity.ParticipationStatus, grade *float64, feedback string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE participations
		SET status = $1, grade = $2, feedback = $3,
		    completed_at = CASE WHEN $1 IN ('completed', 'attended', 'missed') THEN CURRENT_TIMESTAMP ELSE NULL END,
		    points_earned = CASE WHEN $1 = 'completed' OR $1 = 'attended' THEN 
		        (SELECT points FROM activities WHERE id = participations.activity_id) ELSE 0 END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, status, grade, feedback, id)

	if err != nil {
		log.Printf("Database error in UpdateParticipationStatus: %v", err)
		return err
	}
	return nil
}

// Create создаёт новую активность
func (r *ActivityRepository) Create(ctx context.Context, a *activity.Activity) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO activities (
			title, description, type, status,
			start_time, end_time, deadline,
			location, online_link,
			max_participants, points, weight,
			course_id, group_id,
			created_by, created_by_role
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`,
		a.Title, a.Description, a.Type, a.Status,
		a.StartTime, a.EndTime, a.Deadline,
		a.Location, a.OnlineLink,
		a.MaxParticipants, a.Points, a.Weight,
		a.CourseID, a.GroupID,
		a.CreatedBy, a.CreatedByRole,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		log.Printf("Database error in Create: %v", err)
		return err
	}
	return nil
}

// Update обновляет активность
func (r *ActivityRepository) Update(ctx context.Context, a *activity.Activity) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE activities
		SET title = $1, description = $2, type = $3, status = $4,
		    start_time = $5, end_time = $6, deadline = $7,
		    location = $8, online_link = $9,
		    max_participants = $10, points = $11, weight = $12,
		    course_id = $13, group_id = $14,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $15
	`,
		a.Title, a.Description, a.Type, a.Status,
		a.StartTime, a.EndTime, a.Deadline,
		a.Location, a.OnlineLink,
		a.MaxParticipants, a.Points, a.Weight,
		a.CourseID, a.GroupID,
		a.ID,
	)

	if err != nil {
		log.Printf("Database error in Update: %v", err)
		return err
	}
	return nil
}

// Delete удаляет активность
func (r *ActivityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM activities WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Delete: %v", err)
		return err
	}
	return nil
}

// ListActivities получает список активностей с фильтрацией
func (r *ActivityRepository) ListActivities(ctx context.Context, filter activity.ActivityFilter) ([]*activity.Activity, error) {
	query := `
		SELECT id, title, description, type, status,
		       start_time, end_time, deadline,
		       location, online_link,
		       max_participants, current_participants,
		       points, weight,
		       course_id, group_id,
		       created_by, created_by_role,
		       created_at, updated_at
		FROM activities
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	if len(filter.Type) > 0 {
		query += ` AND type = ANY($` + string(rune(argPos)) + `)`
		args = append(args, filter.Type)
		argPos++
	}

	if len(filter.Status) > 0 {
		query += ` AND status = ANY($` + string(rune(argPos)) + `)`
		args = append(args, filter.Status)
		argPos++
	}

	if filter.CourseID != nil {
		query += ` AND course_id = $` + string(rune(argPos))
		args = append(args, *filter.CourseID)
		argPos++
	}

	if filter.GroupID != nil {
		query += ` AND group_id = $` + string(rune(argPos))
		args = append(args, *filter.GroupID)
		argPos++
	}

	if filter.FromDate != nil {
		query += ` AND (start_time >= $` + string(rune(argPos)) + ` OR deadline >= $` + string(rune(argPos)) + `)`
		args = append(args, *filter.FromDate)
		argPos++
	}

	if filter.ToDate != nil {
		query += ` AND (start_time <= $` + string(rune(argPos)) + ` OR deadline <= $` + string(rune(argPos)) + `)`
		args = append(args, *filter.ToDate)
		argPos++
	}

	if filter.Search != "" {
		query += ` AND (title ILIKE $` + string(rune(argPos)) + ` OR description ILIKE $` + string(rune(argPos)) + `)`
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	query += ` ORDER BY start_time NULLS LAST, created_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Database error in ListActivities: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

// GetActivityParticipants получает всех участников активности
func (r *ActivityRepository) GetActivityParticipants(ctx context.Context, activityID uuid.UUID) ([]*activity.Participation, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.activity_id, p.user_id, p.status,
		       p.grade, p.feedback, p.points_earned,
		       p.enrolled_at, p.completed_at,
		       p.created_at, p.updated_at,
		       u.first_name, u.last_name, u.email
		FROM participations p
		JOIN users u ON p.user_id = u.id
		WHERE p.activity_id = $1
		ORDER BY p.enrolled_at
	`, activityID)

	if err != nil {
		log.Printf("Database error in GetActivityParticipants: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanParticipations(rows)
}

// MarkAttendance отмечает посещение
func (r *ActivityRepository) MarkAttendance(ctx context.Context, participationID uuid.UUID, attended bool) error {
	status := activity.ParticipationStatusMissed
	if attended {
		status = activity.ParticipationStatusAttended
	}
	return r.UpdateParticipationStatus(ctx, participationID, status, nil, "")
}

// SetGrade выставляет оценку за активность
func (r *ActivityRepository) SetGrade(ctx context.Context, participationID uuid.UUID, grade float64, feedback string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE participations
		SET grade = $1, feedback = $2,
		    status = 'completed',
		    completed_at = CURRENT_TIMESTAMP,
		    points_earned = (SELECT points FROM activities WHERE id = participations.activity_id),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, grade, feedback, participationID)

	if err != nil {
		log.Printf("Database error in SetGrade: %v", err)
		return err
	}
	return nil
}

// Вспомогательная функция для сканирования активностей
func (r *ActivityRepository) scanActivities(rows pgx.Rows) ([]*activity.Activity, error) {
	var activities []*activity.Activity
	var isEnrolled bool
	var enrollmentStatus *string

	for rows.Next() {
		a := &activity.Activity{}
		var courseID, groupID *uuid.UUID

		// Проверяем, есть ли дополнительные поля (is_enrolled, enrollment_status)
		fieldDescriptions := rows.FieldDescriptions()
		if len(fieldDescriptions) > 20 {
			err := rows.Scan(
				&a.ID, &a.Title, &a.Description, &a.Type, &a.Status,
				&a.StartTime, &a.EndTime, &a.Deadline,
				&a.Location, &a.OnlineLink,
				&a.MaxParticipants, &a.CurrentParticipants,
				&a.Points, &a.Weight,
				&courseID, &groupID,
				&a.CreatedBy, &a.CreatedByRole,
				&a.CreatedAt, &a.UpdatedAt,
				&isEnrolled, &enrollmentStatus,
			)
			if err != nil {
				log.Printf("Error scanning activity row with enrollment: %v", err)
				continue
			}
		} else {
			err := rows.Scan(
				&a.ID, &a.Title, &a.Description, &a.Type, &a.Status,
				&a.StartTime, &a.EndTime, &a.Deadline,
				&a.Location, &a.OnlineLink,
				&a.MaxParticipants, &a.CurrentParticipants,
				&a.Points, &a.Weight,
				&courseID, &groupID,
				&a.CreatedBy, &a.CreatedByRole,
				&a.CreatedAt, &a.UpdatedAt,
			)
			if err != nil {
				log.Printf("Error scanning activity row: %v", err)
				continue
			}
		}

		a.CourseID = courseID
		a.GroupID = groupID
		activities = append(activities, a)
	}

	return activities, nil
}

// Вспомогательная функция для сканирования участий
func (r *ActivityRepository) scanParticipations(rows pgx.Rows) ([]*activity.Participation, error) {
	var participations []*activity.Participation

	for rows.Next() {
		p := &activity.Participation{}
		var activityTitle, activityType, location *string
		var startTime *time.Time

		// Проверяем, есть ли дополнительные поля из JOIN с activities
		fieldDescriptions := rows.FieldDescriptions()
		if len(fieldDescriptions) > 11 {
			err := rows.Scan(
				&p.ID, &p.ActivityID, &p.UserID, &p.Status,
				&p.Grade, &p.Feedback, &p.PointsEarned,
				&p.EnrolledAt, &p.CompletedAt,
				&p.CreatedAt, &p.UpdatedAt,
				&activityTitle, &activityType, &startTime, &location,
			)
			if err != nil {
				log.Printf("Error scanning participation row with activity: %v", err)
				continue
			}
		} else {
			err := rows.Scan(
				&p.ID, &p.ActivityID, &p.UserID, &p.Status,
				&p.Grade, &p.Feedback, &p.PointsEarned,
				&p.EnrolledAt, &p.CompletedAt,
				&p.CreatedAt, &p.UpdatedAt,
			)
			if err != nil {
				log.Printf("Error scanning participation row: %v", err)
				continue
			}
		}

		participations = append(participations, p)
	}

	return participations, nil
}

// GetByCreator получает активности, созданные пользователем
func (r *ActivityRepository) GetByCreator(ctx context.Context, creatorID uuid.UUID) ([]*activity.Activity, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, type, status,
		       start_time, end_time, deadline,
		       location, online_link,
		       max_participants, current_participants,
		       points, weight,
		       course_id, group_id,
		       created_by, created_by_role,
		       created_at, updated_at
		FROM activities
		WHERE created_by = $1
		ORDER BY created_at DESC
	`, creatorID)

	if err != nil {
		log.Printf("Database error in GetByCreator: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}
