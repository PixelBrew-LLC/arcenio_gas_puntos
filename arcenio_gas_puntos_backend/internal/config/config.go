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
	// Determinar el modo de ejecución primero
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release" // Por defecto es release
	}

	// Solo cargar el archivo .env si NO estamos en modo release
	// En release (producción/.exe), usar variables de entorno del sistema
	if ginMode != "release" {
		// Intentar cargar desde el directorio del ejecutable
		exePath, err := os.Executable()
		if err == nil {
			exeDir := filepath.Dir(exePath)
			envPath := filepath.Join(exeDir, ".env")
			_ = godotenv.Load(envPath)
		} else {
			// Fallback: intentar cargar desde el directorio actual
			_ = godotenv.Load()
		}
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "3000"),
		GinMode:    ginMode,

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "arcenio_gas_db"),

		JWTSecret: getEnv("JWT_SECRET", ""),
		Timezone:  getEnv("TZ", "America/Santo_Domingo"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
