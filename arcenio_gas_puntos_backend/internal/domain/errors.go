package domain

import "errors"

// Errores de dominio — independientes de la capa de infraestructura
var (
	// Auth
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
	ErrUnauthorized       = errors.New("unauthorized")

	// Users
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserInactive      = errors.New("user is inactive")
	ErrInvalidPIN        = errors.New("pin must contain only numbers")

	// Clients
	ErrClientNotFound      = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client with this cedula already exists")
	ErrInvalidCedula       = errors.New("invalid cedula format or checksum")

	// Transactions
	ErrInsufficientPoints = errors.New("insufficient points for redemption")
	ErrBelowMinGallons    = errors.New("gallons below minimum required")
	ErrBelowMinRedeem     = errors.New("balance below minimum required for redemption")
	ErrInvalidAmount      = errors.New("amount must be greater than zero")

	// Settings
	ErrSettingNotFound = errors.New("setting not found")
)
