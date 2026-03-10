package usecase

import (
	"context"
	"errors"
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/utils"

	"github.com/google/uuid"
)

type userManagementUsecase struct {
	userRepo domain.UserRepository
	roleRepo domain.RoleRepository
}

func NewUserManagementUsecase(ur domain.UserRepository, rr domain.RoleRepository) domain.UserManagementUsecase {
	return &userManagementUsecase{
		userRepo: ur,
		roleRepo: rr,
	}
}

func (u *userManagementUsecase) CreateUser(ctx context.Context, user *domain.User, createdByUserID string) (*domain.User, error) {
	// 1. Validar que el PIN sea numérico
	if !utils.IsValidPIN(user.Password) {
		return nil, domain.ErrInvalidPIN
	}

	// 2. Verificar cédula única
	existing, err := u.userRepo.GetByCedula(ctx, user.Cedula)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// 3. Verificar username único
	existingByUsername, err := u.userRepo.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// 4. Hash del PIN
	hashedPIN, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// 5. Asignar campos
	createdBy, _ := uuid.Parse(createdByUserID)
	user.ID = uuid.New()
	user.Password = hashedPIN
	user.IsActive = true
	user.CreatedByUserID = &createdBy
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 6. Guardar
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Limpiar password del resultado
	user.Password = ""
	return user, nil
}

func (u *userManagementUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (u *userManagementUsecase) ListByRole(ctx context.Context, roleID uint) ([]*domain.User, error) {
	users, err := u.userRepo.ListByRole(ctx, roleID)
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

func (u *userManagementUsecase) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// 1. Obtener usuario existente
	existing, err := u.userRepo.GetByID(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	// 2. Actualizar campos permitidos
	existing.Nombres = user.Nombres
	existing.Apellidos = user.Apellidos
	existing.Telefono = user.Telefono
	existing.Direccion = user.Direccion
	existing.UpdatedAt = time.Now()

	// Si se envía un nuevo PIN, validar y hashear
	if user.Password != "" {
		if !utils.IsValidPIN(user.Password) {
			return nil, domain.ErrInvalidPIN
		}
		hashedPIN, err := utils.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		existing.Password = hashedPIN
	}

	if err := u.userRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	existing.Password = ""
	return existing, nil
}

func (u *userManagementUsecase) ToggleActive(ctx context.Context, id string, active bool) error {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.IsActive = active
	user.UpdatedAt = time.Now()
	return u.userRepo.Update(ctx, user)
}
