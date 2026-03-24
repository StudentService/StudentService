package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"backend/internal/domain/calendar"
	"backend/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CalendarRepository struct{}

// GetStudentEvents получает все события, доступные студенту
func (r *CalendarRepository) GetStudentEvents(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*calendar.Event, error) {
	query := `
		SELECT DISTINCT e.id, e.title, e.description, e.type, 
		       e.start_time, e.end_time, e.all_day, 
		       e.location, e.online_link,
		       e.course_id, e.group_id, e.user_id,
		       e.created_by, e.created_by_role, e.created_at, e.updated_at,
		       c.name as course_name,
		       g.name as group_name,
		       u.first_name as creator_first_name,
		       u.last_name as creator_last_name
		FROM events e
		LEFT JOIN courses c ON e.course_id = c.id
		LEFT JOIN groups g ON e.group_id = g.id
		LEFT JOIN users u ON e.created_by = u.id
		LEFT JOIN users student ON student.id = $1
		WHERE e.start_time >= $2 AND e.start_time <= $3
		  AND (
		      e.group_id IS NULL OR 
		      e.group_id = student.group_id OR
		      e.user_id = $1 OR
		      e.course_id IN (
		          SELECT course_id FROM groups WHERE id = student.group_id
		      )
		  )
		ORDER BY e.start_time
	`

	rows, err := db.Pool.Query(ctx, query, userID, from, to)
	if err != nil {
		log.Printf("Database error in GetStudentEvents: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

func (r *CalendarRepository) GetGroupEvents(ctx context.Context, groupID uuid.UUID, from, to time.Time) ([]*calendar.Event, error) {
	query := `
		SELECT e.id, e.title, e.description, e.type, 
		       e.start_time, e.end_time, e.all_day, 
		       e.location, e.online_link,
		       e.course_id, e.group_id, e.user_id,
		       e.created_by, e.created_by_role, e.created_at, e.updated_at,
		       c.name as course_name,
		       g.name as group_name,
		       u.first_name as creator_first_name,
		       u.last_name as creator_last_name
		FROM events e
		LEFT JOIN courses c ON e.course_id = c.id
		LEFT JOIN groups g ON e.group_id = g.id
		LEFT JOIN users u ON e.created_by = u.id
		WHERE e.start_time >= $1 AND e.start_time <= $2
		  AND e.group_id = $3
		ORDER BY e.start_time
	`

	rows, err := db.Pool.Query(ctx, query, from, to, groupID)
	if err != nil {
		log.Printf("Database error in GetGroupEvents: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

func (r *CalendarRepository) GetCourseEvents(ctx context.Context, courseID uuid.UUID, from, to time.Time) ([]*calendar.Event, error) {
	query := `
		SELECT e.id, e.title, e.description, e.type, 
		       e.start_time, e.end_time, e.all_day, 
		       e.location, e.online_link,
		       e.course_id, e.group_id, e.user_id,
		       e.created_by, e.created_by_role, e.created_at, e.updated_at,
		       c.name as course_name,
		       g.name as group_name,
		       u.first_name as creator_first_name,
		       u.last_name as creator_last_name
		FROM events e
		LEFT JOIN courses c ON e.course_id = c.id
		LEFT JOIN groups g ON e.group_id = g.id
		LEFT JOIN users u ON e.created_by = u.id
		WHERE e.start_time >= $1 AND e.start_time <= $2
		  AND e.course_id = $3
		ORDER BY e.start_time
	`

	rows, err := db.Pool.Query(ctx, query, from, to, courseID)
	if err != nil {
		log.Printf("Database error in GetCourseEvents: %v", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

func (r *CalendarRepository) GetByID(ctx context.Context, id uuid.UUID) (*calendar.Event, error) {
	e := &calendar.Event{}
	var courseID, groupID, userID *uuid.UUID

	err := db.Pool.QueryRow(ctx, `
		SELECT id, title, description, type, 
		       start_time, end_time, all_day, 
		       location, online_link,
		       course_id, group_id, user_id,
		       created_by, created_by_role, created_at, updated_at
		FROM events WHERE id = $1
	`, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.Type,
		&e.StartTime, &e.EndTime, &e.AllDay,
		&e.Location, &e.OnlineLink,
		&courseID, &groupID, &userID,
		&e.CreatedBy, &e.CreatedByRole, &e.CreatedAt, &e.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByID: %v", err)
		return nil, err
	}

	e.CourseID = courseID
	e.GroupID = groupID
	e.UserID = userID

	return e, nil
}

func (r *CalendarRepository) Create(ctx context.Context, e *calendar.Event) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO events (
			title, description, type, 
			start_time, end_time, all_day,
			location, online_link,
			course_id, group_id, user_id,
			created_by, created_by_role
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`, e.Title, e.Description, e.Type,
		e.StartTime, e.EndTime, e.AllDay,
		e.Location, e.OnlineLink,
		e.CourseID, e.GroupID, e.UserID,
		e.CreatedBy, e.CreatedByRole,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		log.Printf("Database error in Create: %v", err)
		return err
	}
	return nil
}

func (r *CalendarRepository) Update(ctx context.Context, e *calendar.Event) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE events
		SET title = $1, description = $2, type = $3,
		    start_time = $4, end_time = $5, all_day = $6,
		    location = $7, online_link = $8,
		    course_id = $9, group_id = $10, user_id = $11,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
	`, e.Title, e.Description, e.Type,
		e.StartTime, e.EndTime, e.AllDay,
		e.Location, e.OnlineLink,
		e.CourseID, e.GroupID, e.UserID,
		e.ID)

	if err != nil {
		log.Printf("Database error in Update: %v", err)
		return err
	}
	return nil
}

func (r *CalendarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Delete: %v", err)
		return err
	}
	return nil
}

// Вспомогательная функция для сканирования строк
func (r *CalendarRepository) scanEvents(rows pgx.Rows) ([]*calendar.Event, error) {
	var events []*calendar.Event

	for rows.Next() {
		e := &calendar.Event{}
		var courseID, groupID, userID *uuid.UUID
		var courseName, groupName, creatorFirstName, creatorLastName sql.NullString // 👈 используем sql.NullString

		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Type,
			&e.StartTime, &e.EndTime, &e.AllDay,
			&e.Location, &e.OnlineLink,
			&courseID, &groupID, &userID,
			&e.CreatedBy, &e.CreatedByRole, &e.CreatedAt, &e.UpdatedAt,
			&courseName, &groupName,
			&creatorFirstName, &creatorLastName,
		)
		if err != nil {
			log.Printf("Error scanning event row: %v", err)
			continue
		}

		e.CourseID = courseID
		e.GroupID = groupID
		e.UserID = userID

		// Сохраняем дополнительные поля (можно добавить в структуру Event при необходимости)
		_ = courseName.String
		_ = groupName.String
		_ = creatorFirstName.String + " " + creatorLastName.String

		events = append(events, e)
	}

	return events, nil
}
