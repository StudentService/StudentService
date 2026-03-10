package infrastructure

import (
	"context"

	"backend/internal/domain/user"
	"backend/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct{}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	u := &user.User{}
	err := db.Pool.QueryRow(ctx, `SELECT id, username, email, role, first_name, last_name FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.FirstName, &u.LastName)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	// TODO: реализовать позже
	return nil
	// или return errors.New("not implemented")
}
