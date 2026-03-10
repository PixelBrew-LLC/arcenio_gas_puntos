package handler

import (
	"net/http"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api/dto"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsUsecase domain.SettingsUsecase
}

func NewSettingsHandler(su domain.SettingsUsecase) *SettingsHandler {
	return &SettingsHandler{settingsUsecase: su}
}

func (h *SettingsHandler) GetAll(c *gin.Context) {
	settings, err := h.settingsUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener configuraciones"})
		return
	}

	responses := make([]dto.SettingResponse, len(settings))
	for i, s := range settings {
		responses[i] = dto.SettingResponse{
			Key:   s.Key,
			Value: s.Value,
		}
	}

	c.JSON(http.StatusOK, responses)
}

func (h *SettingsHandler) Update(c *gin.Context) {
	var req dto.UpdateSettingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsUsecase.Update(c.Request.Context(), req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar configuración"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configuración actualizada"})
}
