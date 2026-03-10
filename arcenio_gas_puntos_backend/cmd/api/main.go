package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/config"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/database"

	"github.com/gin-gonic/gin"
)

// fatalWait imprime el error y espera Enter antes de cerrar (útil en Windows .exe)
func fatalWait(format string, args ...interface{}) {
	log.Printf(format, args...)
	fmt.Println("\nPresiona Enter para cerrar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}

func main() {
	// 1. Cargar configuración (.env)
	cfg := config.LoadConfig()

	// 2. Conectar a la base de datos
	db, err := database.NewConnection(cfg)
	if err != nil {
		fatalWait("❌ Error al conectar con la base de datos: %v", err)
	}
	log.Println("✅ Conexión a PostgreSQL establecida")

	// 3. Ejecutar migraciones
	if err := database.Migrate(db); err != nil {
		fatalWait("❌ Error al ejecutar migraciones: %v", err)
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
		fatalWait("❌ Error al iniciar el servidor: %v", err)
	}
}
