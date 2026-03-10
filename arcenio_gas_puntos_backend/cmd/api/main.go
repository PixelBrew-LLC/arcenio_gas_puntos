package main

import (
	"log"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/config"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Cargar configuración (.env)
	cfg := config.LoadConfig()

	// 2. Conectar a la base de datos
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Error al conectar con la base de datos: %v", err)
	}
	log.Println("✅ Conexión a PostgreSQL establecida")

	// 3. Ejecutar migraciones
	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ Error al ejecutar migraciones: %v", err)
	}
	log.Println("✅ Migraciones ejecutadas")

	// 4. Sembrar datos iniciales
	database.SeedAll(db)

	// 5. Configurar modo de Gin
	gin.SetMode(cfg.GinMode)

	// 6. Inicializar el router con inyección de dependencias
	r := api.NewRouter(db, cfg)

	// 7. Iniciar el servidor
	log.Printf("🚀 Servidor iniciando en puerto %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor: %v", err)
	}
}
