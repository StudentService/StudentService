package repositories

import (
	"context"
	"log"

	"backend/internal/domain/group"
	"backend/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GroupRepository struct{}

func (r *GroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*group.Group, error) {
	g := &group.Group{}
	var holderID *uuid.UUID

	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, course_id, semester_id, holder_id, created_at, updated_at 
		FROM groups WHERE id = $1
	`, id).Scan(
		&g.ID, &g.Name, &g.CourseID, &g.SemesterID, &holderID, &g.CreatedAt, &g.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in Group.GetByID: %v", err)
		return nil, err
	}

	g.HolderID = holderID
	return g, nil
}

func (r *GroupRepository) GetByStudent(ctx context.Context, studentID uuid.UUID) (*group.Group, error) {
	g := &group.Group{}
	var holderID *uuid.UUID

	// Сначала получаем group_id пользователя, затем группу
	err := db.Pool.QueryRow(ctx, `
		SELECT g.id, g.name, g.course_id, g.semester_id, g.holder_id, g.created_at, g.updated_at
		FROM groups g
		JOIN users u ON u.group_id = g.id
		WHERE u.id = $1
	`, studentID).Scan(
		&g.ID, &g.Name, &g.CourseID, &g.SemesterID, &holderID, &g.CreatedAt, &g.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in Group.GetByStudent: %v", err)
		return nil, err
	}

	g.HolderID = holderID
	return g, nil
}

func (r *GroupRepository) GetAll(ctx context.Context) ([]*group.Group, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, course_id, semester_id, holder_id, created_at, updated_at 
		FROM groups ORDER BY name
	`)
	if err != nil {
		log.Printf("Database error in Group.GetAll: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*group.Group
	for rows.Next() {
		g := &group.Group{}
		var holderID *uuid.UUID
		err := rows.Scan(
			&g.ID, &g.Name, &g.CourseID, &g.SemesterID, &holderID, &g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning group row: %v", err)
			continue
		}
		g.HolderID = holderID
		groups = append(groups, g)
	}

	return groups, nil
}

func (r *GroupRepository) Create(ctx context.Context, g *group.Group) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO groups (name, course_id, semester_id, holder_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, g.Name, g.CourseID, g.SemesterID, g.HolderID).Scan(
		&g.ID, &g.CreatedAt, &g.UpdatedAt,
	)

	if err != nil {
		log.Printf("Database error in Group.Create: %v", err)
		return err
	}
	return nil
}

func (r *GroupRepository) Update(ctx context.Context, g *group.Group) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE groups
		SET name = $1, course_id = $2, semester_id = $3, holder_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`, g.Name, g.CourseID, g.SemesterID, g.HolderID, g.ID)

	if err != nil {
		log.Printf("Database error in Group.Update: %v", err)
		return err
	}
	return nil
}

func (r *GroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM groups WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Group.Delete: %v", err)
		return err
	}
	return nil
}

// Дополнительный метод для получения групп по курсу
func (r *GroupRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*group.Group, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, course_id, semester_id, holder_id, created_at, updated_at 
		FROM groups WHERE course_id = $1
	`, courseID)
	if err != nil {
		log.Printf("Database error in Group.GetByCourseID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*group.Group
	for rows.Next() {
		g := &group.Group{}
		var holderID *uuid.UUID
		err := rows.Scan(
			&g.ID, &g.Name, &g.CourseID, &g.SemesterID, &holderID, &g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning group row: %v", err)
			continue
		}
		g.HolderID = holderID
		groups = append(groups, g)
	}

	return groups, nil
}
