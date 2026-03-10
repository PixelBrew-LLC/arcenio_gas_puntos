package handler

import (
	"errors"
	"net/http"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api/dto"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(au domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: au}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authUsecase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		statusCode, errorMessage := mapAuthError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	roleName := ""
	if result.User.Role != nil {
		roleName = result.User.Role.Name
	}

	res := dto.LoginResponse{
		User: dto.UserData{
			ID:        result.User.ID.String(),
			Nombres:   result.User.Nombres,
			Apellidos: result.User.Apellidos,
			Username:  result.User.Username,
			Role:      roleName,
		},
		AccessToken: result.AccessToken,
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	user, err := h.authUsecase.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		statusCode, errorMessage := mapAuthError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	res := dto.MeResponse{
		ID:        user.ID.String(),
		Nombres:   user.Nombres,
		Apellidos: user.Apellidos,
		Cedula:    user.Cedula,
		Telefono:  user.Telefono,
		Username:  user.Username,
		Role:      roleName,
		IsActive:  user.IsActive,
	}

	c.JSON(http.StatusOK, res)
}

func mapAuthError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "credenciales inválidas"
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "usuario no encontrado"
	case errors.Is(err, domain.ErrUserInactive):
		return http.StatusForbidden, "usuario inactivo"
	case errors.Is(err, domain.ErrInvalidToken):
		return http.StatusUnauthorized, "token inválido"
	case errors.Is(err, domain.ErrExpiredToken):
		return http.StatusUnauthorized, "token expirado"
	default:
		return http.StatusInternalServerError, "error interno del servidor"
	}
}
