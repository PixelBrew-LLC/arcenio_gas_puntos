package database

import (
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/repository/postgres"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&postgres.RoleModel{},
		&postgres.UserModel{},
		&postgres.ClientModel{},
		&postgres.PointsLedgerModel{},
		&postgres.SettingModel{},
	)
}
