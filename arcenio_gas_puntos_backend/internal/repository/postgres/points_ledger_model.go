package postgres

import (
	"time"

	"github.com/google/uuid"
)

// PointsLedgerModel es el modelo de base de datos para el ledger de puntos
type PointsLedgerModel struct {
	ID                uuid.UUID   `gorm:"type:uuid;primary_key"`
	ClientID          uuid.UUID   `gorm:"type:uuid;not null;index"`
	Client            ClientModel `gorm:"foreignKey:ClientID"`
	Points            float64     `gorm:"type:decimal(10,2);not null"`
	TransactionType   string      `gorm:"type:varchar(20);not null"` // earn, redeem, expire
	GallonsAmount     float64     `gorm:"type:decimal(10,2);default:0"`
	ProcessedByUserID uuid.UUID   `gorm:"type:uuid;not null"`
	ProcessedByUser   UserModel   `gorm:"foreignKey:ProcessedByUserID"`
	CreatedAt         time.Time
	ExpiresAt         *time.Time `gorm:"index"`
}

func (PointsLedgerModel) TableName() string {
	return "points_ledger"
}
