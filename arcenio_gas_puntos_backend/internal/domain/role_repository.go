package domain

import "context"

// RoleRepository define las operaciones de persistencia de roles
type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*Role, error)
	GetAll(ctx context.Context) ([]*Role, error)
}
