package domain

import "context"

// UserManagementUsecase define las operaciones CRUD para gestión de usuarios (Bomberos/Admins)
type UserManagementUsecase interface {
	CreateUser(ctx context.Context, user *User, createdByUserID string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	ListByRole(ctx context.Context, roleID uint) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	ToggleActive(ctx context.Context, id string, active bool) error
}
