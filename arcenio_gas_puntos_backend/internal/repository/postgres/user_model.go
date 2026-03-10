package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserModel es el modelo de base de datos para usuarios
type UserModel struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key"`
	Nombres         string     `gorm:"type:varchar(100);not null"`
	Apellidos       string     `gorm:"type:varchar(100);not null"`
	Cedula          string     `gorm:"type:varchar(30);uniqueIndex;not null"`
	Telefono        string     `gorm:"type:varchar(20);not null"`
	Direccion       string     `gorm:"type:varchar(255);not null"`
	Username        string     `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password        string     `gorm:"type:varchar(255);not null"`
	RoleID          uint       `gorm:"not null"`
	Role            RoleModel  `gorm:"foreignKey:RoleID"`
	IsActive        bool       `gorm:"default:true"`
	CreatedByUserID *uuid.UUID `gorm:"type:uuid"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
	return "users"
}
