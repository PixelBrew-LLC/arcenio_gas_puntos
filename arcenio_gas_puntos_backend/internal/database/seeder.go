package database

import (
	"log"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/repository/postgres"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	seedRoles(db)
	seedSettings(db)
	seedSuperAdmin(db)
}

func seedRoles(db *gorm.DB) {
	roles := []string{domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleBombero}

	var count int64
	db.Model(&postgres.RoleModel{}).Count(&count)
	if count > 0 {
		return
	}

	for _, r := range roles {
		db.Create(&postgres.RoleModel{Name: r})
	}
	log.Println("✅ Seeded Roles: SuperAdmin, Admin, Bombero")
}

func seedSettings(db *gorm.DB) {
	defaults := map[string]string{
		domain.SettingPointsPerGallon:    "5",
		domain.SettingMinGallons:         "1",
		domain.SettingMinRedeemPoints:    "500",
		domain.SettingPointsExpiryMonths: "12",
	}

	var count int64
	db.Model(&postgres.SettingModel{}).Count(&count)
	if count > 0 {
		return
	}

	for key, value := range defaults {
		db.Create(&postgres.SettingModel{Key: key, Value: value})
	}
	log.Println("✅ Seeded Default Settings")
}

func seedSuperAdmin(db *gorm.DB) {
	// Verificar si ya existe un SuperAdmin
	var role postgres.RoleModel
	if err := db.Where("name = ?", domain.RoleSuperAdmin).First(&role).Error; err != nil {
		log.Println("⚠️  SuperAdmin role not found, skipping SuperAdmin seed")
		return
	}

	var count int64
	db.Model(&postgres.UserModel{}).Where("role_id = ?", role.ID).Count(&count)
	if count > 0 {
		return
	}

	hashedPIN, err := utils.HashPassword("1234")
	if err != nil {
		log.Printf("⚠️  Error hashing SuperAdmin PIN: %v", err)
		return
	}

	superAdmin := &postgres.UserModel{
		ID:        uuid.New(),
		Nombres:   "Super",
		Apellidos: "Administrador",
		Cedula:    "0000000000",
		Telefono:  "0000000000",
		Direccion: "Oficina",
		Username:  "admin",
		Password:  hashedPIN,
		RoleID:    role.ID,
		IsActive:  true,
	}

	if err := db.Create(superAdmin).Error; err != nil {
		log.Printf("⚠️  Error creating SuperAdmin: %v", err)
		return
	}
	log.Println("✅ Seeded SuperAdmin (username: admin, PIN: 1234)")
}
