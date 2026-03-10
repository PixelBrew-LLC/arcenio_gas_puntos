package dto

import "time"

// EarnPointsRequest es el DTO para acumular puntos
type EarnPointsRequest struct {
	ClientID string  `json:"client_id" binding:"required"`
	Gallons  float64 `json:"gallons" binding:"required,gt=0"`
}

// RedeemPointsRequest es el DTO para canjear puntos (siempre canjea todos)
type RedeemPointsRequest struct {
	ClientID string `json:"client_id" binding:"required"`
}

// TransactionResponse es el DTO de respuesta para una transacción
type TransactionResponse struct {
	ID              string    `json:"id"`
	ClientID        string    `json:"client_id"`
	Points          float64   `json:"points"`
	TransactionType string    `json:"transaction_type"`
	GallonsAmount   float64   `json:"gallons_amount"`
	ProcessedBy     string    `json:"processed_by"`
	ProcessedByName string    `json:"processed_by_name"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       *string   `json:"expires_at,omitempty"`
}

// EarnResponse es el DTO de respuesta para una acumulación
type EarnResponse struct {
	PointsEarned float64             `json:"points_earned"`
	NewBalance   float64             `json:"new_balance"`
	Transaction  TransactionResponse `json:"transaction"`
	Client       ClientResponse      `json:"client"`
}

// RedeemResponse es el DTO de respuesta para un canje
type RedeemResponse struct {
	PointsRedeemed float64             `json:"points_redeemed"`
	NewBalance     float64             `json:"new_balance"`
	Transaction    TransactionResponse `json:"transaction"`
	Client         ClientResponse      `json:"client"`
}

// BalanceResponse es el DTO de respuesta para el saldo
type BalanceResponse struct {
	ClientID  string  `json:"client_id"`
	Balance   float64 `json:"balance"`
	MinRedeem float64 `json:"min_redeem"`
}

// DashboardTopClient es un cliente destacado en el dashboard
type DashboardTopClient struct {
	Nombres   string  `json:"nombres"`
	Apellidos string  `json:"apellidos"`
	Cedula    string  `json:"cedula"`
	Points    float64 `json:"points"`
}

// DashboardRecentTransaction es una transacción reciente en el dashboard
type DashboardRecentTransaction struct {
	TransactionType string  `json:"transaction_type"`
	ClientName      string  `json:"client_name"`
	Points          float64 `json:"points"`
	GallonsAmount   float64 `json:"gallons_amount"`
	CreatedAt       string  `json:"created_at"`
}

// DashboardResponse es el DTO de respuesta para el dashboard
type DashboardResponse struct {
	TotalGallons        float64                      `json:"total_gallons"`
	TotalPointsEarned   float64                      `json:"total_points_earned"`
	TotalPointsRedeemed float64                      `json:"total_points_redeemed"`
	TotalTransactions   int64                        `json:"total_transactions"`
	TotalClients        int64                        `json:"total_clients"`
	Month               int                          `json:"month"`
	Year                int                          `json:"year"`
	TopClients          []DashboardTopClient         `json:"top_clients"`
	RecentTransactions  []DashboardRecentTransaction `json:"recent_transactions"`
}
