package user

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error
}
