package postgres

import (
	"context"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type clientRepo struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) domain.ClientRepository {
	return &clientRepo{db: db}
}

func (r *clientRepo) toDomain(model *ClientModel) *domain.Client {
	if model == nil {
		return nil
	}
	return &domain.Client{
		ID:                 model.ID,
		Nombres:            model.Nombres,
		Apellidos:          model.Apellidos,
		Cedula:             model.Cedula,
		Direccion:          model.Direccion,
		Telefono:           model.Telefono,
		RegisteredByUserID: model.RegisteredByUserID,
		CreatedAt:          model.CreatedAt,
	}
}

func (r *clientRepo) toModel(client *domain.Client) *ClientModel {
	if client == nil {
		return nil
	}
	return &ClientModel{
		ID:                 client.ID,
		Nombres:            client.Nombres,
		Apellidos:          client.Apellidos,
		Cedula:             client.Cedula,
		Direccion:          client.Direccion,
		Telefono:           client.Telefono,
		RegisteredByUserID: client.RegisteredByUserID,
		CreatedAt:          client.CreatedAt,
	}
}

func (r *clientRepo) Create(ctx context.Context, client *domain.Client) error {
	model := r.toModel(client)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *clientRepo) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	clientID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrClientNotFound
	}

	var model ClientModel
	err = r.db.WithContext(ctx).First(&model, clientID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *clientRepo) GetByCedula(ctx context.Context, cedula string) (*domain.Client, error) {
	var model ClientModel
	err := r.db.WithContext(ctx).Where("cedula = ?", cedula).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *clientRepo) List(ctx context.Context) ([]*domain.Client, error) {
	// Usamos query builder para sumar los puntos no expirados.
	// La expiración cuenta solo para transacciones 'earn' donde expires_at > NOW() o es NULL.
	// Redenciones ('redeem') y expiraciones ('expire') siempre cuentan hacia el saldo.
	var results []struct {
		ClientModel
		PointsBalance float64
	}

	err := r.db.WithContext(ctx).
		Table("clients").
		Select("clients.*, COALESCE(SUM(points_ledger.points), 0) as points_balance").
		Joins("LEFT JOIN points_ledger ON clients.id = points_ledger.client_id AND (points_ledger.expires_at IS NULL OR points_ledger.expires_at > CURRENT_TIMESTAMP)").
		Group("clients.id").
		Order("clients.created_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	clients := make([]*domain.Client, len(results))
	for i := range results {
		domainClient := r.toDomain(&results[i].ClientModel)
		domainClient.PointsBalance = results[i].PointsBalance
		clients[i] = domainClient
	}
	return clients, nil
}

func (r *clientRepo) Update(ctx context.Context, client *domain.Client) error {
	model := r.toModel(client)
	return r.db.WithContext(ctx).Save(model).Error
}
