package middleware

import (
	"net/http"
	"strings"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireAuth verifica que el request tenga un JWT válido
func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extraer el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorización requerido"})
			c.Abort()
			return
		}

		// 2. El formato debe ser: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de autorización inválido"})
			c.Abort()
			return
		}

		token := parts[1]

		// 3. Validar el JWT
		claims, err := utils.ValidateAccessToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		// 4. Inyectar datos del usuario en el contexto
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Set("roleName", claims.RoleName)

		c.Next()
	}
}
