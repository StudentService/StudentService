package user

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, u *User) error
	// ... другие методы
}
