package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api/dto"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase domain.TransactionUsecase
	tz                 *time.Location
}

func NewTransactionHandler(tu domain.TransactionUsecase, tz *time.Location) *TransactionHandler {
	return &TransactionHandler{transactionUsecase: tu, tz: tz}
}

func (h *TransactionHandler) EarnPoints(c *gin.Context) {
	var req dto.EarnPointsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	result, err := h.transactionUsecase.EarnPoints(
		c.Request.Context(),
		req.ClientID,
		req.Gallons,
		userID.(string),
	)
	if err != nil {
		statusCode, errorMessage := mapTransactionError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, dto.EarnResponse{
		PointsEarned: result.PointsEarned,
		NewBalance:   result.NewBalance,
		Transaction:  h.toTransactionResponse(result.Transaction),
		Client:       toClientResponse(result.Client),
	})
}

func (h *TransactionHandler) RedeemPoints(c *gin.Context) {
	var req dto.RedeemPointsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	result, err := h.transactionUsecase.RedeemPoints(
		c.Request.Context(),
		req.ClientID,
		userID.(string),
	)
	if err != nil {
		statusCode, errorMessage := mapTransactionError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, dto.RedeemResponse{
		PointsRedeemed: result.PointsRedeemed,
		NewBalance:     result.NewBalance,
		Transaction:    h.toTransactionResponse(result.Transaction),
		Client:         toClientResponse(result.Client),
	})
}

func (h *TransactionHandler) GetBalance(c *gin.Context) {
	clientID := c.Param("clientId")

	balance, err := h.transactionUsecase.GetClientBalance(c.Request.Context(), clientID)
	if err != nil {
		statusCode, errorMessage := mapTransactionError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	minRedeem, _ := h.transactionUsecase.GetMinRedeemPoints(c.Request.Context())

	c.JSON(http.StatusOK, dto.BalanceResponse{
		ClientID:  clientID,
		Balance:   balance,
		MinRedeem: minRedeem,
	})
}

func (h *TransactionHandler) GetHistory(c *gin.Context) {
	filter := domain.TransactionFilter{}

	if clientID := c.Query("client_id"); clientID != "" {
		filter.ClientID = &clientID
	}
	if userID := c.Query("user_id"); userID != "" {
		filter.UserID = &userID
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter.DateFrom = &t
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			endOfDay := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.DateTo = &endOfDay
		}
	}

	entries, err := h.transactionUsecase.GetTransactionHistory(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener historial"})
		return
	}

	responses := make([]dto.TransactionResponse, len(entries))
	for i, entry := range entries {
		responses[i] = h.toTransactionResponse(entry)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *TransactionHandler) GetClientHistory(c *gin.Context) {
	clientID := c.Param("clientId")

	filter := domain.TransactionFilter{
		ClientID: &clientID,
	}

	entries, err := h.transactionUsecase.GetTransactionHistory(c.Request.Context(), filter)
	if err != nil {
		statusCode, errorMessage := mapTransactionError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	responses := make([]dto.TransactionResponse, len(entries))
	for i, entry := range entries {
		responses[i] = h.toTransactionResponse(entry)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *TransactionHandler) GetDashboard(c *gin.Context) {
	now := time.Now().In(h.tz)
	month := int(now.Month())
	year := now.Year()

	if m := c.Query("month"); m != "" {
		if parsed, err := strconv.Atoi(m); err == nil {
			month = parsed
		}
	}
	if y := c.Query("year"); y != "" {
		if parsed, err := strconv.Atoi(y); err == nil {
			year = parsed
		}
	}

	stats, err := h.transactionUsecase.GetDashboardStats(c.Request.Context(), month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener estadísticas"})
		return
	}

	// Map top clients
	topClients := make([]dto.DashboardTopClient, len(stats.TopClients))
	for i, tc := range stats.TopClients {
		topClients[i] = dto.DashboardTopClient{
			Nombres:   tc.Nombres,
			Apellidos: tc.Apellidos,
			Cedula:    tc.Cedula,
			Points:    tc.Points,
		}
	}

	// Map recent transactions
	recentTx := make([]dto.DashboardRecentTransaction, len(stats.RecentTransactions))
	for i, rt := range stats.RecentTransactions {
		recentTx[i] = dto.DashboardRecentTransaction{
			TransactionType: rt.TransactionType,
			ClientName:      rt.ClientName,
			Points:          rt.Points,
			GallonsAmount:   rt.GallonsAmount,
			CreatedAt:       rt.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.DashboardResponse{
		TotalGallons:        stats.TotalGallons,
		TotalPointsEarned:   stats.TotalPointsEarned,
		TotalPointsRedeemed: stats.TotalPointsRedeemed,
		TotalTransactions:   stats.TotalTransactions,
		TotalClients:        stats.TotalClients,
		Month:               month,
		Year:                year,
		TopClients:          topClients,
		RecentTransactions:  recentTx,
	})
}

func (h *TransactionHandler) toTransactionResponse(entry *domain.PointsLedger) dto.TransactionResponse {
	resp := dto.TransactionResponse{
		ID:              entry.ID.String(),
		ClientID:        entry.ClientID.String(),
		Points:          entry.Points,
		TransactionType: string(entry.TransactionType),
		GallonsAmount:   entry.GallonsAmount,
		ProcessedBy:     entry.ProcessedByUserID.String(),
		ProcessedByName: entry.ProcessedByName,
		CreatedAt:       entry.CreatedAt.In(h.tz),
	}
	if entry.ExpiresAt != nil {
		formatted := entry.ExpiresAt.In(h.tz).Format(time.RFC3339)
		resp.ExpiresAt = &formatted
	}
	return resp
}

func mapTransactionError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrClientNotFound):
		return http.StatusNotFound, "cliente no encontrado"
	case errors.Is(err, domain.ErrBelowMinGallons):
		return http.StatusBadRequest, "los galones están por debajo del mínimo requerido"
	case errors.Is(err, domain.ErrInsufficientPoints):
		return http.StatusBadRequest, "puntos insuficientes para el canje"
	case errors.Is(err, domain.ErrBelowMinRedeem):
		return http.StatusBadRequest, "el saldo está por debajo del mínimo para canjear"
	case errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusBadRequest, "la cantidad debe ser mayor a cero"
	default:
		return http.StatusInternalServerError, "error interno del servidor"
	}
}
