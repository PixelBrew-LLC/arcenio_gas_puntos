package postgres

import (
	"context"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"gorm.io/gorm"
)

type roleRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &roleRepo{db: db}
}

func (r *roleRepo) toDomain(model *RoleModel) *domain.Role {
	if model == nil {
		return nil
	}
	return &domain.Role{
		ID:   model.ID,
		Name: model.Name,
	}
}

func (r *roleRepo) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	var model RoleModel
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *roleRepo) GetAll(ctx context.Context) ([]*domain.Role, error) {
	var models []RoleModel
	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	roles := make([]*domain.Role, len(models))
	for i := range models {
		roles[i] = r.toDomain(&models[i])
	}
	return roles, nil
}
