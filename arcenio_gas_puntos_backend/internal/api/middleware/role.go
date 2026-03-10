package middleware

import (
	"net/http"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/gin-gonic/gin"
)

// RequireRole verifica que el usuario tenga uno de los roles permitidos
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName, exists := c.Get("roleName")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
			c.Abort()
			return
		}

		userRole := roleName.(string)

		// SuperAdmin siempre tiene acceso a todo
		if userRole == domain.RoleSuperAdmin {
			c.Next()
			return
		}

		// Verificar si el rol del usuario está en la lista de permitidos
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "no tiene permisos para esta operación"})
		c.Abort()
	}
}
