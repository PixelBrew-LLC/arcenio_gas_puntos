package domain

import "context"

// TransactionUsecase define las operaciones de negocio para transacciones de puntos
type TransactionUsecase interface {
	EarnPoints(ctx context.Context, clientID string, gallons float64, processedByUserID string) (*EarnResult, error)
	RedeemPoints(ctx context.Context, clientID string, processedByUserID string) (*RedeemResult, error)
	GetClientBalance(ctx context.Context, clientID string) (float64, error)
	GetMinRedeemPoints(ctx context.Context) (float64, error)
	GetTransactionHistory(ctx context.Context, filter TransactionFilter) ([]*PointsLedger, error)
	GetDashboardStats(ctx context.Context, month, year int) (*DashboardStats, error)
}
