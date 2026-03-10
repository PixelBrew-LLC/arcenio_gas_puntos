package domain

import "context"

// LoginResult contiene el token y usuario después del login
type LoginResult struct {
	User        *User
	AccessToken string
}

// AuthUsecase define las operaciones de autenticación
type AuthUsecase interface {
	Login(ctx context.Context, username, pin string) (*LoginResult, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
}
