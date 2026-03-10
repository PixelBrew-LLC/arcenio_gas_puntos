package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role representa un rol del sistema
type Role struct {
	ID   uint
	Name string
}

// Constantes de roles
const (
	RoleSuperAdmin = "SuperAdmin"
	RoleAdmin      = "Admin"
	RoleBombero    = "Bombero"
)

// User representa un usuario del sistema (Bombero, Admin o SuperAdmin)
type User struct {
	ID              uuid.UUID
	Nombres         string
	Apellidos       string
	Cedula          string
	Telefono        string
	Direccion       string
	Username        string
	Password        string // Hash bcrypt
	RoleID          uint
	Role            *Role
	IsActive        bool
	CreatedByUserID *uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
