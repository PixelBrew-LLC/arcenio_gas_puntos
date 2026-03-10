package usecase

import (
	"context"
	"errors"
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/google/uuid"
)

type clientUsecase struct {
	clientRepo domain.ClientRepository
}

func NewClientUsecase(cr domain.ClientRepository) domain.ClientUsecase {
	return &clientUsecase{clientRepo: cr}
}

func (u *clientUsecase) Create(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	if err := client.ValidateCedula(); err != nil {
		return nil, err
	}

	// 1. Verificar si ya existe un cliente con esa cédula
	existing, err := u.clientRepo.GetByCedula(ctx, client.Cedula)
	if err != nil && !errors.Is(err, domain.ErrClientNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrClientAlreadyExists
	}

	// 2. Asignar ID y timestamp
	client.ID = uuid.New()
	client.CreatedAt = time.Now()

	// 3. Guardar
	if err := u.clientRepo.Create(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

func (u *clientUsecase) GetByCedula(ctx context.Context, cedula string) (*domain.Client, error) {
	return u.clientRepo.GetByCedula(ctx, cedula)
}

func (u *clientUsecase) List(ctx context.Context) ([]*domain.Client, error) {
	return u.clientRepo.List(ctx)
}

func (u *clientUsecase) Update(ctx context.Context, id string, updates *domain.Client) (*domain.Client, error) {
	existing, err := u.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Nombres = updates.Nombres
	existing.Apellidos = updates.Apellidos
	if updates.Direccion != nil {
		existing.Direccion = updates.Direccion
	}
	if updates.Telefono != nil {
		existing.Telefono = updates.Telefono
	}

	if err := u.clientRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}
