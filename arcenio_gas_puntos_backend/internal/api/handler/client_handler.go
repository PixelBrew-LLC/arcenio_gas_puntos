package handler

import (
	"errors"
	"net/http"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api/dto"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientUsecase domain.ClientUsecase
}

func NewClientHandler(cu domain.ClientUsecase) *ClientHandler {
	return &ClientHandler{clientUsecase: cu}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.CreateClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	client := &domain.Client{
		Nombres:   req.Nombres,
		Apellidos: req.Apellidos,
		Cedula:    req.Cedula,
		Direccion: req.Direccion,
		Telefono:  req.Telefono,
	}

	// Parsear el userID
	if uid, ok := userID.(string); ok {
		parsedID, err := parseUUID(uid)
		if err == nil {
			client.RegisteredByUserID = parsedID
		}
	}

	result, err := h.clientUsecase.Create(c.Request.Context(), client)
	if err != nil {
		statusCode, errorMessage := mapClientError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, toClientResponse(result))
}

func (h *ClientHandler) GetByCedula(c *gin.Context) {
	cedula := c.Param("cedula")

	client, err := h.clientUsecase.GetByCedula(c.Request.Context(), cedula)
	if err != nil {
		statusCode, errorMessage := mapClientError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, toClientResponse(client))
}

func (h *ClientHandler) List(c *gin.Context) {
	clients, err := h.clientUsecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener clientes"})
		return
	}

	responses := make([]dto.ClientResponse, len(clients))
	for i, client := range clients {
		responses[i] = toClientResponse(client)
	}

	c.JSON(http.StatusOK, responses)
}

func toClientResponse(client *domain.Client) dto.ClientResponse {
	return dto.ClientResponse{
		ID:            client.ID.String(),
		Nombres:       client.Nombres,
		Apellidos:     client.Apellidos,
		Cedula:        client.Cedula,
		Direccion:     client.Direccion,
		Telefono:      client.Telefono,
		PointsBalance: client.PointsBalance,
		CreatedAt:     client.CreatedAt,
	}
}

func mapClientError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrClientNotFound):
		return http.StatusNotFound, "cliente no encontrado"
	case errors.Is(err, domain.ErrClientAlreadyExists):
		return http.StatusConflict, "ya existe un cliente con esta cédula"
	case errors.Is(err, domain.ErrInvalidCedula):
		return http.StatusBadRequest, "la cédula ingresada no es válida"
	default:
		return http.StatusInternalServerError, "error interno del servidor"
	}
}

func (h *ClientHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := &domain.Client{
		Nombres:   req.Nombres,
		Apellidos: req.Apellidos,
		Direccion: req.Direccion,
		Telefono:  req.Telefono,
	}

	result, err := h.clientUsecase.Update(c.Request.Context(), id, updates)
	if err != nil {
		statusCode, errorMessage := mapClientError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, toClientResponse(result))
}
