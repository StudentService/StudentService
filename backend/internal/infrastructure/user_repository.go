package infrastructure

import (
	"context"
	"log"

	"backend/internal/domain/user"
	"backend/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct{}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	u := &user.User{}
	err := db.Pool.QueryRow(ctx, `
		SELECT id, username, email, password_hash, role, first_name, last_name, created_at, updated_at 
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByID: %v", err)
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	err := db.Pool.QueryRow(ctx, `
        INSERT INTO users (username, email, password_hash, role, first_name, last_name)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `, u.Username, u.Email, u.PasswordHash, u.Role, u.FirstName, u.LastName).Scan(
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
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $7
    `, u.Username, u.Email, u.PasswordHash, u.Role, u.FirstName, u.LastName, u.ID)

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
	err := db.Pool.QueryRow(ctx, `
		SELECT id, username, email, password_hash, role, first_name, last_name, created_at, updated_at 
		FROM users WHERE email = $1
	`, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		log.Printf("Database error in GetByEmail: %v", err)
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	u := &user.User{}
	err := db.Pool.QueryRow(ctx, `
		SELECT id, username, email, password_hash, role, first_name, last_name, created_at, updated_at 
		FROM users WHERE username = $1
	`, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Printf("Database error in GetByUsername: %v", err)
		return nil, err
	}
	return u, nil
}
