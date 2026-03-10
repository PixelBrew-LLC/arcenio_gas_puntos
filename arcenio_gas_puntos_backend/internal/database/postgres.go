package database

import (
	"fmt"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
