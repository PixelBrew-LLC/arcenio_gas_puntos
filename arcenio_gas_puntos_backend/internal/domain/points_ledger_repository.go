package domain

import (
	"context"
	"time"
)

// TransactionFilter contiene los filtros para consultar transacciones
type TransactionFilter struct {
	ClientID *string
	UserID   *string
	DateFrom *time.Time
	DateTo   *time.Time
}

// DashboardStats contiene las estadísticas del dashboard
type DashboardStats struct {
	TotalGallons        float64
	TotalPointsEarned   float64
	TotalPointsRedeemed float64
	TotalTransactions   int64
	TotalClients        int64
	TopClients          []TopClient
	RecentTransactions  []RecentTransaction
}

// TopClient representa un cliente destacado por puntos
type TopClient struct {
	Nombres   string
	Apellidos string
	Cedula    string
	Points    float64
}

// RecentTransaction representa una transacción reciente
type RecentTransaction struct {
	TransactionType string
	ClientName      string
	Points          float64
	GallonsAmount   float64
	CreatedAt       string
}

// PointsLedgerRepository define las operaciones de persistencia del ledger
type PointsLedgerRepository interface {
	Create(ctx context.Context, entry *PointsLedger) error
	GetBalance(ctx context.Context, clientID string) (float64, error)
	GetByClientID(ctx context.Context, clientID string) ([]*PointsLedger, error)
	ListFiltered(ctx context.Context, filter TransactionFilter) ([]*PointsLedger, error)
	ExpireOldPoints(ctx context.Context) (int64, error)
	GetDashboardStats(ctx context.Context, month, year int) (*DashboardStats, error)
}
