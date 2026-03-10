package dto

import "time"

// CreateClientRequest es el DTO para crear un cliente
type CreateClientRequest struct {
	Nombres   string  `json:"nombres" binding:"required"`
	Apellidos string  `json:"apellidos" binding:"required"`
	Cedula    string  `json:"cedula" binding:"required"`
	Direccion *string `json:"direccion"`
	Telefono  *string `json:"telefono"`
}

// UpdateClientRequest es el DTO para actualizar un cliente
type UpdateClientRequest struct {
	Nombres   string  `json:"nombres" binding:"required"`
	Apellidos string  `json:"apellidos" binding:"required"`
	Direccion *string `json:"direccion"`
	Telefono  *string `json:"telefono"`
}

// ClientResponse es el DTO de respuesta para un cliente
type ClientResponse struct {
	ID            string    `json:"id"`
	Nombres       string    `json:"nombres"`
	Apellidos     string    `json:"apellidos"`
	Cedula        string    `json:"cedula"`
	Direccion     *string   `json:"direccion"`
	Telefono      *string   `json:"telefono"`
	PointsBalance float64   `json:"points_balance"`
	CreatedAt     time.Time `json:"created_at"`
}
