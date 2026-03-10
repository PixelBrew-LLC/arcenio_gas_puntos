package domain

import "context"

// SettingsUsecase define las operaciones de negocio para configuraciones globales
type SettingsUsecase interface {
	GetAll(ctx context.Context) ([]*Setting, error)
	Update(ctx context.Context, key, value string) error
}
