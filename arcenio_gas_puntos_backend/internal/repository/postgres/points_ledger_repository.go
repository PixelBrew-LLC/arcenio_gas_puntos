package postgres

import (
	"context"
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pointsLedgerRepo struct {
	db *gorm.DB
}

func NewPointsLedgerRepository(db *gorm.DB) domain.PointsLedgerRepository {
	return &pointsLedgerRepo{db: db}
}

func (r *pointsLedgerRepo) toDomain(model *PointsLedgerModel) *domain.PointsLedger {
	if model == nil {
		return nil
	}

	processedByName := ""
	if model.ProcessedByUser.ID != uuid.Nil {
		processedByName = model.ProcessedByUser.Nombres + " " + model.ProcessedByUser.Apellidos
	}

	return &domain.PointsLedger{
		ID:                model.ID,
		ClientID:          model.ClientID,
		Points:            model.Points,
		TransactionType:   domain.TransactionType(model.TransactionType),
		GallonsAmount:     model.GallonsAmount,
		ProcessedByUserID: model.ProcessedByUserID,
		ProcessedByName:   processedByName,
		CreatedAt:         model.CreatedAt,
		ExpiresAt:         model.ExpiresAt,
	}
}

func (r *pointsLedgerRepo) toModel(entry *domain.PointsLedger) *PointsLedgerModel {
	if entry == nil {
		return nil
	}
	return &PointsLedgerModel{
		ID:                entry.ID,
		ClientID:          entry.ClientID,
		Points:            entry.Points,
		TransactionType:   string(entry.TransactionType),
		GallonsAmount:     entry.GallonsAmount,
		ProcessedByUserID: entry.ProcessedByUserID,
		CreatedAt:         entry.CreatedAt,
		ExpiresAt:         entry.ExpiresAt,
	}
}

func (r *pointsLedgerRepo) Create(ctx context.Context, entry *domain.PointsLedger) error {
	model := r.toModel(entry)
	return r.db.WithContext(ctx).Create(model).Error
}

// GetBalance calcula el saldo disponible: suma de puntos donde no ha expirado
func (r *pointsLedgerRepo) GetBalance(ctx context.Context, clientID string) (float64, error) {
	cID, err := uuid.Parse(clientID)
	if err != nil {
		return 0, domain.ErrClientNotFound
	}

	var balance float64
	err = r.db.WithContext(ctx).
		Model(&PointsLedgerModel{}).
		Where("client_id = ? AND (expires_at IS NULL OR expires_at > ?)", cID, time.Now()).
		Select("COALESCE(SUM(points), 0)").
		Scan(&balance).Error

	return balance, err
}

func (r *pointsLedgerRepo) GetByClientID(ctx context.Context, clientID string) ([]*domain.PointsLedger, error) {
	cID, err := uuid.Parse(clientID)
	if err != nil {
		return nil, domain.ErrClientNotFound
	}

	var models []PointsLedgerModel
	err = r.db.WithContext(ctx).
		Preload("ProcessedByUser").
		Where("client_id = ?", cID).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entries := make([]*domain.PointsLedger, len(models))
	for i := range models {
		entries[i] = r.toDomain(&models[i])
	}
	return entries, nil
}

func (r *pointsLedgerRepo) ListFiltered(ctx context.Context, filter domain.TransactionFilter) ([]*domain.PointsLedger, error) {
	query := r.db.WithContext(ctx).Model(&PointsLedgerModel{}).Order("created_at DESC")

	if filter.ClientID != nil {
		cID, err := uuid.Parse(*filter.ClientID)
		if err == nil {
			query = query.Where("client_id = ?", cID)
		}
	}
	if filter.UserID != nil {
		uID, err := uuid.Parse(*filter.UserID)
		if err == nil {
			query = query.Where("processed_by_user_id = ?", uID)
		}
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	var models []PointsLedgerModel
	err := query.Preload("ProcessedByUser").Find(&models).Error
	if err != nil {
		return nil, err
	}

	entries := make([]*domain.PointsLedger, len(models))
	for i := range models {
		entries[i] = r.toDomain(&models[i])
	}
	return entries, nil
}

// ExpireOldPoints marca como expirados los puntos cuya fecha de expiración ya pasó
func (r *pointsLedgerRepo) ExpireOldPoints(ctx context.Context) (int64, error) {
	// Buscar entries de tipo "earn" que hayan expirado y cuyo saldo aún no se haya restado
	now := time.Now()

	// Obtener todas las entries de earn que ya expiraron
	var expiredEntries []PointsLedgerModel
	err := r.db.WithContext(ctx).
		Where("transaction_type = ? AND expires_at IS NOT NULL AND expires_at <= ?", string(domain.TransactionTypeEarn), now).
		Find(&expiredEntries).Error
	if err != nil {
		return 0, err
	}

	// Se manejan las expiraciones a través del filtro de balance (expires_at > now)
	// No necesitamos crear entries adicionales de "expire" ya que el balance
	// solo suma entries cuyo expires_at > now
	return int64(len(expiredEntries)), nil
}

func (r *pointsLedgerRepo) GetDashboardStats(ctx context.Context, month, year int) (*domain.DashboardStats, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	stats := &domain.DashboardStats{}

	// Total galones despachados (earn transactions)
	r.db.WithContext(ctx).
		Model(&PointsLedgerModel{}).
		Where("transaction_type = ? AND created_at >= ? AND created_at < ?", string(domain.TransactionTypeEarn), startDate, endDate).
		Select("COALESCE(SUM(gallons_amount), 0)").
		Scan(&stats.TotalGallons)

	// Total puntos otorgados
	r.db.WithContext(ctx).
		Model(&PointsLedgerModel{}).
		Where("transaction_type = ? AND created_at >= ? AND created_at < ?", string(domain.TransactionTypeEarn), startDate, endDate).
		Select("COALESCE(SUM(points), 0)").
		Scan(&stats.TotalPointsEarned)

	// Total puntos canjeados (valor absoluto)
	r.db.WithContext(ctx).
		Model(&PointsLedgerModel{}).
		Where("transaction_type = ? AND created_at >= ? AND created_at < ?", string(domain.TransactionTypeRedeem), startDate, endDate).
		Select("COALESCE(ABS(SUM(points)), 0)").
		Scan(&stats.TotalPointsRedeemed)

	// Total transacciones del mes
	r.db.WithContext(ctx).
		Model(&PointsLedgerModel{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Count(&stats.TotalTransactions)

	// Total clientes registrados (global, no solo del mes)
	r.db.WithContext(ctx).
		Model(&ClientModel{}).
		Count(&stats.TotalClients)

	// Top 5 clientes por puntos acumulados en el mes
	type topClientRow struct {
		Nombres   string
		Apellidos string
		Cedula    string
		Points    float64
	}
	var topRows []topClientRow
	r.db.WithContext(ctx).
		Table("points_ledger").
		Select("clients.nombres, clients.apellidos, clients.cedula, SUM(points_ledger.points) as points").
		Joins("JOIN clients ON clients.id = points_ledger.client_id").
		Where("points_ledger.transaction_type = ? AND points_ledger.created_at >= ? AND points_ledger.created_at < ?",
			string(domain.TransactionTypeEarn), startDate, endDate).
		Group("clients.id, clients.nombres, clients.apellidos, clients.cedula").
		Order("points DESC").
		Limit(5).
		Scan(&topRows)

	for _, row := range topRows {
		stats.TopClients = append(stats.TopClients, domain.TopClient{
			Nombres:   row.Nombres,
			Apellidos: row.Apellidos,
			Cedula:    row.Cedula,
			Points:    row.Points,
		})
	}

	// 5 transacciones más recientes
	type recentRow struct {
		TransactionType string
		ClientName      string
		Points          float64
		GallonsAmount   float64
		CreatedAt       time.Time
	}
	var recentRows []recentRow
	r.db.WithContext(ctx).
		Table("points_ledger").
		Select("points_ledger.transaction_type, CONCAT(clients.nombres, ' ', clients.apellidos) as client_name, points_ledger.points, points_ledger.gallons_amount, points_ledger.created_at").
		Joins("JOIN clients ON clients.id = points_ledger.client_id").
		Order("points_ledger.created_at DESC").
		Limit(5).
		Scan(&recentRows)

	for _, row := range recentRows {
		stats.RecentTransactions = append(stats.RecentTransactions, domain.RecentTransaction{
			TransactionType: row.TransactionType,
			ClientName:      row.ClientName,
			Points:          row.Points,
			GallonsAmount:   row.GallonsAmount,
			CreatedAt:       row.CreatedAt.Format(time.RFC3339),
		})
	}

	return stats, nil
}
