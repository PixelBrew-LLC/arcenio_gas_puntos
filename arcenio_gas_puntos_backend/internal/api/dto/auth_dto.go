package dto

// LoginRequest es el DTO para login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse es el DTO de respuesta con token
type LoginResponse struct {
	User        UserData `json:"user"`
	AccessToken string   `json:"access_token"`
}

// UserData contiene los datos básicos del usuario autenticado
type UserData struct {
	ID        string `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

// MeResponse es el DTO de respuesta para /me
type MeResponse struct {
	ID        string `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Cedula    string `json:"cedula"`
	Telefono  string `json:"telefono"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
}
