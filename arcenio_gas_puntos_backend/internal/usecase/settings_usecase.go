package usecase

import (
	"context"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"
)

type settingsUsecase struct {
	settingsRepo domain.SettingsRepository
}

func NewSettingsUsecase(sr domain.SettingsRepository) domain.SettingsUsecase {
	return &settingsUsecase{settingsRepo: sr}
}

func (u *settingsUsecase) GetAll(ctx context.Context) ([]*domain.Setting, error) {
	return u.settingsRepo.GetAll(ctx)
}

func (u *settingsUsecase) Update(ctx context.Context, key, value string) error {
	return u.settingsRepo.Upsert(ctx, key, value)
}
