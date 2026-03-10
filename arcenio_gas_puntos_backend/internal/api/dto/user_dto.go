package dto

// CreateUserRequest es el DTO para crear un usuario (Bombero o Admin)
type CreateUserRequest struct {
	Nombres   string `json:"nombres" binding:"required"`
	Apellidos string `json:"apellidos" binding:"required"`
	Cedula    string `json:"cedula" binding:"required"`
	Telefono  string `json:"telefono" binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=4"` // PIN numérico
	RoleID    uint   `json:"role_id" binding:"required"`
}

// UpdateUserRequest es el DTO para actualizar un usuario
type UpdateUserRequest struct {
	Nombres   string `json:"nombres" binding:"required"`
	Apellidos string `json:"apellidos" binding:"required"`
	Telefono  string `json:"telefono" binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
	Password  string `json:"password,omitempty"` // Nuevo PIN (opcional)
}

// ToggleActiveRequest es el DTO para activar/desactivar un usuario
type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

// UserResponse es el DTO de respuesta para un usuario
type UserResponse struct {
	ID        string `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Cedula    string `json:"cedula"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	Username  string `json:"username"`
	RoleID    uint   `json:"role_id"`
	RoleName  string `json:"role_name"`
	IsActive  bool   `json:"is_active"`
}

// SettingResponse es el DTO de respuesta para una configuración
type SettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// UpdateSettingRequest es el DTO para actualizar una configuración
type UpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
