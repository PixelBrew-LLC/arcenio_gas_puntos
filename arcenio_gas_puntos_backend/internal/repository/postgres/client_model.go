package postgres

import (
	"time"

	"github.com/google/uuid"
)

// ClientModel es el modelo de base de datos para clientes
type ClientModel struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key"`
	Nombres            string    `gorm:"type:varchar(100);not null"`
	Apellidos          string    `gorm:"type:varchar(100);not null"`
	Cedula             string    `gorm:"type:varchar(30);uniqueIndex;not null"`
	Direccion          *string   `gorm:"type:varchar(255)"`
	Telefono           *string   `gorm:"type:varchar(20)"`
	RegisteredByUserID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt          time.Time
}

func (ClientModel) TableName() string {
	return "clients"
}
