package repositories

import (
	"context"
	"log"

	"backend/internal/domain/user"
	"backend/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct{}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	u := &user.User{}
	var groupID *uuid.UUID

	err := db.Pool.QueryRow(ctx, `
		SELECT id, username, email, password_hash, role, first_name, last_name, group_id, created_at, updated_at 
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &groupID, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByID: %v", err)
		return nil, err
	}

	u.GroupID = groupID
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	err := db.Pool.QueryRow(ctx, `
        INSERT INTO users (username, email, password_hash, role, first_name, last_name, group_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `, u.Username, u.Email, u.PasswordHash, u.Role, u.FirstName, u.LastName, u.GroupID).Scan(
		&u.ID, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		log.Printf("Database error in Create: %v", err)
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	_, err := db.Pool.Exec(ctx, `
        UPDATE users
        SET username = $1, email = $2, password_hash = $3, 
            role = $4, first_name = $5, last_name = $6, 
            group_id = $7, updated_at = CURRENT_TIMESTAMP
        WHERE id = $8
    `, u.Username, u.Email, u.PasswordHash, u.Role, u.FirstName, u.LastName, u.GroupID, u.ID)

	if err != nil {
		log.Printf("Database error in Update: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Delete: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u := &user.User{}
	var groupID *uuid.UUID

	err := db.Pool.QueryRow(ctx, `
		SELECT id, username, email, password_hash, role, first_name, last_name, group_id, created_at, updated_at 
		FROM users WHERE email = $1
	`, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &groupID, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Printf("Database error in GetByEmail: %v", err)
		return nil, err
	}

	u.GroupID = groupID
	return u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	u := &user.User{}
	var groupID *uuid.UUID

	err := db.Pool.QueryRow(ctx, `
		SELECT id, username, email, password_hash, role, first_name, last_name, group_id, created_at, updated_at 
		FROM users WHERE username = $1
	`, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &groupID, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Printf("Database error in GetByUsername: %v", err)
		return nil, err
	}

	u.GroupID = groupID
	return u, nil
}

// Дополнительный метод для получения пользователей по группе
func (r *UserRepository) GetByGroupID(ctx context.Context, groupID uuid.UUID) ([]*user.User, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, username, email, password_hash, role, first_name, last_name, group_id, created_at, updated_at 
		FROM users WHERE group_id = $1
	`, groupID)
	if err != nil {
		log.Printf("Database error in GetByGroupID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := &user.User{}
		var gid *uuid.UUID
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
			&u.FirstName, &u.LastName, &gid, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}
		u.GroupID = gid
		users = append(users, u)
	}

	return users, nil
}
