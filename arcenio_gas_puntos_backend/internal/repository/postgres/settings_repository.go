package postgres

import (
	"context"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type settingsRepo struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) domain.SettingsRepository {
	return &settingsRepo{db: db}
}

func (r *settingsRepo) toDomain(model *SettingModel) *domain.Setting {
	if model == nil {
		return nil
	}
	return &domain.Setting{
		Key:   model.Key,
		Value: model.Value,
	}
}

func (r *settingsRepo) Get(ctx context.Context, key string) (*domain.Setting, error) {
	var model SettingModel
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSettingNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *settingsRepo) GetAll(ctx context.Context) ([]*domain.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	settings := make([]*domain.Setting, len(models))
	for i := range models {
		settings[i] = r.toDomain(&models[i])
	}
	return settings, nil
}

func (r *settingsRepo) Upsert(ctx context.Context, key, value string) error {
	model := SettingModel{Key: key, Value: value}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&model).Error
}
