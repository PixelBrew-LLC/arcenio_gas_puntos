package domain

import (
	"time"

	"github.com/google/uuid"
)

// Client representa un cliente del programa de fidelización
type Client struct {
	ID                 uuid.UUID
	Nombres            string
	Apellidos          string
	Cedula             string
	Direccion          *string
	Telefono           *string
	RegisteredByUserID uuid.UUID
	CreatedAt          time.Time
	PointsBalance      float64 // Computed field
}

// ValidateCedula valida que la cédula tenga 11 dígitos y cumpla el algoritmo de Módulo 10.
func (c *Client) ValidateCedula() error {
	if len(c.Cedula) != 11 {
		return ErrInvalidCedula
	}

	for _, ch := range c.Cedula {
		if ch < '0' || ch > '9' {
			return ErrInvalidCedula
		}
	}

	multipliers := []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2}
	suma := 0

	for i := 0; i < 10; i++ {
		val := int(c.Cedula[i] - '0')
		res := val * multipliers[i]

		if res >= 10 {
			suma += (res / 10) + (res % 10)
		} else {
			suma += res
		}
	}

	checkDigit := int(c.Cedula[10] - '0')

	decenaSuperior := suma - (suma % 10) + 10
	total := decenaSuperior - suma

	// Ajuste estándar del cálculo de cédula dominicana: si el total da 10, debe validarse como 0
	if total == 10 {
		total = 0
	}

	if total != checkDigit {
		return ErrInvalidCedula
	}

	return nil
}
