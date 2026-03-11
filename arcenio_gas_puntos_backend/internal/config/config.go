package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string

	// Gin
	GinMode string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret string

	// Timezone
	Timezone string
}

func LoadConfig() *Config {
	// Obtener el directorio del ejecutable
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		envPath := filepath.Join(exeDir, ".env")
		_ = godotenv.Load(envPath)
	} else {
		// Fallback: intentar cargar desde el directorio actual
		_ = godotenv.Load()
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "3000"),
		GinMode:    getEnv("GIN_MODE", "release"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "manuel"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "arcenio_gas_db"),

		JWTSecret: getEnv("JWT_SECRET", "default_secret"),
		Timezone:  getEnv("TZ", "America/Tegucigalpa"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
