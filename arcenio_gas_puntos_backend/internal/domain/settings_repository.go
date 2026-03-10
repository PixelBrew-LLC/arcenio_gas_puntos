package domain

import "context"

// SettingsRepository define las operaciones de persistencia de configuraciones
type SettingsRepository interface {
	Get(ctx context.Context, key string) (*Setting, error)
	GetAll(ctx context.Context) ([]*Setting, error)
	Upsert(ctx context.Context, key, value string) error
}
