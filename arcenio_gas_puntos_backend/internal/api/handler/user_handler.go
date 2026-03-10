package handler

import (
	"errors"
	"net/http"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/api/dto"
	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userMgmtUsecase domain.UserManagementUsecase
}

func NewUserHandler(umu domain.UserManagementUsecase) *UserHandler {
	return &UserHandler{userMgmtUsecase: umu}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdByID, _ := c.Get("userID")

	user := &domain.User{
		Nombres:   req.Nombres,
		Apellidos: req.Apellidos,
		Cedula:    req.Cedula,
		Telefono:  req.Telefono,
		Direccion: req.Direccion,
		Username:  req.Username,
		Password:  req.Password,
		RoleID:    req.RoleID,
	}

	result, err := h.userMgmtUsecase.CreateUser(c.Request.Context(), user, createdByID.(string))
	if err != nil {
		statusCode, errorMessage := mapUserError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(result))
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userMgmtUsecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		statusCode, errorMessage := mapUserError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) ListBomberos(c *gin.Context) {
	h.listByRoleName(c, domain.RoleBombero)
}

func (h *UserHandler) ListAdmins(c *gin.Context) {
	h.listByRoleName(c, domain.RoleAdmin)
}

func (h *UserHandler) listByRoleName(c *gin.Context, roleName string) {
	// Necesitamos el roleID. Lo extraemos del query param o lo hardcodeamos.
	// En la práctica, los roles son: 1=SuperAdmin, 2=Admin, 3=Bombero (seeded en ese orden)
	var roleID uint
	switch roleName {
	case domain.RoleSuperAdmin:
		roleID = 1
	case domain.RoleAdmin:
		roleID = 2
	case domain.RoleBombero:
		roleID = 3
	}

	users, err := h.userMgmtUsecase.ListByRole(c.Request.Context(), roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener usuarios"})
		return
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = toUserResponse(user)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user := &domain.User{
		ID:        parsedID,
		Nombres:   req.Nombres,
		Apellidos: req.Apellidos,
		Telefono:  req.Telefono,
		Direccion: req.Direccion,
		Password:  req.Password,
	}

	result, err := h.userMgmtUsecase.UpdateUser(c.Request.Context(), user)
	if err != nil {
		statusCode, errorMessage := mapUserError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(result))
}

func (h *UserHandler) ToggleActive(c *gin.Context) {
	id := c.Param("id")
	var req dto.ToggleActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userMgmtUsecase.ToggleActive(c.Request.Context(), id, req.IsActive); err != nil {
		statusCode, errorMessage := mapUserError(err)
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado"})
}

func toUserResponse(user *domain.User) dto.UserResponse {
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	return dto.UserResponse{
		ID:        user.ID.String(),
		Nombres:   user.Nombres,
		Apellidos: user.Apellidos,
		Cedula:    user.Cedula,
		Telefono:  user.Telefono,
		Direccion: user.Direccion,
		Username:  user.Username,
		RoleID:    user.RoleID,
		RoleName:  roleName,
		IsActive:  user.IsActive,
	}
}

func mapUserError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "usuario no encontrado"
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, "ya existe un usuario con esta cédula o nombre de usuario"
	case errors.Is(err, domain.ErrInvalidPIN):
		return http.StatusBadRequest, "el PIN debe contener solo números (mínimo 4 dígitos)"
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusForbidden, "no autorizado para esta operación"
	default:
		return http.StatusInternalServerError, "error interno del servidor"
	}
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
