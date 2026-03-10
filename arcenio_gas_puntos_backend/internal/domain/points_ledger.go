package domain

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType define el tipo de transacción en el ledger
type TransactionType string

const (
	TransactionTypeEarn   TransactionType = "earn"
	TransactionTypeRedeem TransactionType = "redeem"
	TransactionTypeExpire TransactionType = "expire"
)

// PointsLedger representa un movimiento de puntos (acumulación, canje o expiración)
type PointsLedger struct {
	ID                uuid.UUID
	ClientID          uuid.UUID
	Points            float64 // Positivo al ganar, negativo al canjear/expirar
	TransactionType   TransactionType
	GallonsAmount     float64
	ProcessedByUserID uuid.UUID
	ProcessedByName   string // populated when Preloaded
	CreatedAt         time.Time
	ExpiresAt         *time.Time // nil para redeem/expire
}

// EarnResult es el resultado de una operación de acumulación
type EarnResult struct {
	PointsEarned float64
	NewBalance   float64
	Transaction  *PointsLedger
	Client       *Client
}

// RedeemResult es el resultado de una operación de canje
type RedeemResult struct {
	PointsRedeemed float64
	NewBalance     float64
	Transaction    *PointsLedger
	Client         *Client
}
