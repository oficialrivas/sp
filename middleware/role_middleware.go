package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role") // Asegúrate de que el rol del usuario esté almacenado en el contexto después de la autenticación JWT

		for _, role := range roles {
			if strings.EqualFold(role, userRole) {
				c.Next()
				return
			}
		}

		// Obtener el área y nombre de la entidad del contexto
		entityArea, areaExists := c.Get("entityArea")
		entityName, nameExists := c.Get("entityName")

		// Crear mensaje de error detallado
		errorMessage := "You don't have permission to access this resource"
		if areaExists && nameExists {
			errorMessage = "You don't have permission to access this resource in area: " + entityArea.(string) + ", entity: " + entityName.(string)
		}

		c.JSON(http.StatusForbidden, gin.H{"error": errorMessage})
		c.Abort()
	}
}
