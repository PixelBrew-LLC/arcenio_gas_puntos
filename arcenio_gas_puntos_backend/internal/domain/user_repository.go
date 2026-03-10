package domain

import "context"

// UserRepository define las operaciones de persistencia de usuarios
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByCedula(ctx context.Context, cedula string) (*User, error)
	ListByRole(ctx context.Context, roleID uint) ([]*User, error)
	Update(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, id string) error
}
