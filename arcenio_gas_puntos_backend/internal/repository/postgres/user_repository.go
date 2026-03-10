package postgres

import (
	"context"
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) toDomain(model *UserModel) *domain.User {
	if model == nil {
		return nil
	}

	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	var role *domain.Role
	if model.Role.ID != 0 {
		role = &domain.Role{
			ID:   model.Role.ID,
			Name: model.Role.Name,
		}
	}

	return &domain.User{
		ID:              model.ID,
		Nombres:         model.Nombres,
		Apellidos:       model.Apellidos,
		Cedula:          model.Cedula,
		Telefono:        model.Telefono,
		Direccion:       model.Direccion,
		Username:        model.Username,
		Password:        model.Password,
		RoleID:          model.RoleID,
		Role:            role,
		IsActive:        model.IsActive,
		CreatedByUserID: model.CreatedByUserID,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}

func (r *userRepo) toModel(user *domain.User) *UserModel {
	if user == nil {
		return nil
	}

	model := &UserModel{
		ID:              user.ID,
		Nombres:         user.Nombres,
		Apellidos:       user.Apellidos,
		Cedula:          user.Cedula,
		Telefono:        user.Telefono,
		Direccion:       user.Direccion,
		Username:        user.Username,
		Password:        user.Password,
		RoleID:          user.RoleID,
		IsActive:        user.IsActive,
		CreatedByUserID: user.CreatedByUserID,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}

	if user.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{
			Time:  *user.DeletedAt,
			Valid: true,
		}
	}

	return model
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	model := r.toModel(user)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	var model UserModel
	err = r.db.WithContext(ctx).Preload("Role").First(&model, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *userRepo) GetByCedula(ctx context.Context, cedula string) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Preload("Role").Where("cedula = ?", cedula).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *userRepo) ListByRole(ctx context.Context, roleID uint) ([]*domain.User, error) {
	var models []UserModel
	err := r.db.WithContext(ctx).Preload("Role").Where("role_id = ?", roleID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i := range models {
		users[i] = r.toDomain(&models[i])
	}
	return users, nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	model := r.toModel(user)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *userRepo) SoftDelete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrUserNotFound
	}
	return r.db.WithContext(ctx).Delete(&UserModel{}, userID).Error
}
