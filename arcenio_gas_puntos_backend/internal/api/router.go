package api

import (
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api/handler"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api/middleware"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/config"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/health"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/repository/postgres"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(middleware.CORS())

	// --- TIMEZONE ---
	tz, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		tz = time.UTC
	}

	// --- INYECCIÓN DE DEPENDENCIAS ---

	// 1. Repositorios
	roleRepo := postgres.NewRoleRepository(db)
	userRepo := postgres.NewUserRepository(db)
	clientRepo := postgres.NewClientRepository(db)
	ledgerRepo := postgres.NewPointsLedgerRepository(db)
	settingsRepo := postgres.NewSettingsRepository(db)

	// Suprimir warning de variable no usada
	_ = roleRepo

	// 2. Casos de Uso
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret, tz)
	clientUsecase := usecase.NewClientUsecase(clientRepo)
	transactionUsecase := usecase.NewTransactionUsecase(ledgerRepo, clientRepo, settingsRepo)
	userMgmtUsecase := usecase.NewUserManagementUsecase(userRepo, roleRepo)
	settingsUsecase := usecase.NewSettingsUsecase(settingsRepo)

	// 3. Handlers
	healthHandler := health.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authUsecase)
	clientHandler := handler.NewClientHandler(clientUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase, tz)
	userHandler := handler.NewUserHandler(userMgmtUsecase)
	settingsHandler := handler.NewSettingsHandler(settingsUsecase)

	// --- RUTAS ---

	// Health check
	r.GET("/health", healthHandler.Check)

	// Auth (público)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middleware.RequireAuth(cfg.JWTSecret), authHandler.GetMe)
	}

	// Clientes (Bombero+)
	clientGroup := r.Group("/clients")
	clientGroup.Use(middleware.RequireAuth(cfg.JWTSecret))
	clientGroup.Use(middleware.RequireRole(domain.RoleBombero, domain.RoleAdmin))
	{
		clientGroup.POST("", clientHandler.Create)
		clientGroup.GET("", clientHandler.List)
		clientGroup.GET("/:cedula", clientHandler.GetByCedula)
		clientGroup.PUT("/:id", clientHandler.Update)
	}

	// Transacciones (Bombero+)
	txGroup := r.Group("/transactions")
	txGroup.Use(middleware.RequireAuth(cfg.JWTSecret))
	txGroup.Use(middleware.RequireRole(domain.RoleBombero, domain.RoleAdmin))
	{
		txGroup.POST("/earn", transactionHandler.EarnPoints)
		txGroup.POST("/redeem", transactionHandler.RedeemPoints)
		txGroup.GET("/balance/:clientId", transactionHandler.GetBalance)
		txGroup.GET("/history/:clientId", transactionHandler.GetClientHistory)
	}

	// Reportes (Admin+)
	reportGroup := r.Group("/reports")
	reportGroup.Use(middleware.RequireAuth(cfg.JWTSecret))
	reportGroup.Use(middleware.RequireRole(domain.RoleAdmin))
	{
		reportGroup.GET("/transactions", transactionHandler.GetHistory)
		reportGroup.GET("/dashboard", transactionHandler.GetDashboard)
	}

	// Gestión de Bomberos (Admin+)
	bomberoGroup := r.Group("/users/bomberos")
	bomberoGroup.Use(middleware.RequireAuth(cfg.JWTSecret))
	bomberoGroup.Use(middleware.RequireRole(domain.RoleAdmin))
	{
		bomberoGroup.POST("", userHandler.Create)
		bomberoGroup.GET("", userHandler.ListBomberos)
		bomberoGroup.GET("/:id", userHandler.GetByID)
		bomberoGroup.PUT("/:id", userHandler.Update)
		bomberoGroup.PATCH("/:id/active", userHandler.ToggleActive)
	}

	// Gestión de Admins (SuperAdmin only)
	adminGroup := r.Group("/users/admins")
	adminGroup.Use(middleware.RequireAuth(cfg.JWTSecret))
	adminGroup.Use(middleware.RequireRole(domain.RoleSuperAdmin))
	{
		adminGroup.POST("", userHandler.Create)
		adminGroup.GET("", userHandler.ListAdmins)
		adminGroup.GET("/:id", userHandler.GetByID)
		adminGroup.PUT("/:id", userHandler.Update)
		adminGroup.PATCH("/:id/active", userHandler.ToggleActive)
	}

	// Configuraciones (Admin+)
	settingsGroup := r.Group("/settings")
	settingsGroup.Use(middleware.RequireAuth(cfg.JWTSecret))
	settingsGroup.Use(middleware.RequireRole(domain.RoleAdmin))
	{
		settingsGroup.GET("", settingsHandler.GetAll)
		settingsGroup.PUT("", settingsHandler.Update)
	}

	return r
}
