package usecase

import (
	"context"
	"errors"
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/utils"

	"github.com/google/uuid"
)

type authUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
	tz        *time.Location
}

func NewAuthUsecase(ur domain.UserRepository, jwtSecret string, tz *time.Location) domain.AuthUsecase {
	return &authUsecase{
		userRepo:  ur,
		jwtSecret: jwtSecret,
		tz:        tz,
	}
}

func (u *authUsecase) Login(ctx context.Context, username, pin string) (*domain.LoginResult, error) {
	// 1. Buscar usuario por username
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// 3. Verificar PIN
	if !utils.CheckPasswordHash(pin, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// 4. Generar JWT con expiración a medianoche del día actual
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Username,
		u.jwtSecret,
		user.RoleID,
		roleName,
		u.tz,
	)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResult{
		User:        user,
		AccessToken: accessToken,
	}, nil
}

func (u *authUsecase) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	_ = id
	return u.userRepo.GetByID(ctx, userID)
}
