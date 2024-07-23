package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	configs "github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
)

func AreaCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			c.Abort()
			return
		}

		var user models.User
		if err := configs.DB.Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		path := c.FullPath()
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported entity type"})
			c.Abort()
			return
		}

		entityType := parts[1]
		entityID := c.Param("id")
		var entityArea string
		var entityName string

		switch entityType {
		case "casos":
			var entity models.Caso
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Caso"
			}

		case "documentos":
			var entity models.Documento
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Documento"
			}

		case "pasaportes":
			var entity models.Pasaporte
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Pasaporte"
			}

		case "personas":
			var entity models.Persona
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Persona"
			}

		case "vehiculos":
			var entity models.Vehiculo
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Vehiculo"
			}

		case "empresas":
			var entity models.Empresa
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Empresa"
			}

		case "direcciones":
			var entity models.Direccion
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Direccion"
			}

		case "visas":
			var entity models.Visa
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "Visa"
			}

		case "iios":
			var entity models.IIO
			if entityID != "" {
				if err := configs.DB.Where("id = ?", entityID).First(&entity).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
					c.Abort()
					return
				}
				entityArea = entity.Area
				entityName = "IIO"
			}

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported entity type"})
			c.Abort()
			return
		}

		// Guardar información de entidad en el contexto
		c.Set("entityArea", entityArea)
		c.Set("entityName", entityName)

		// Permitir acceso si el usuario es admin
		if user.Nivel == "admin" {
			c.Next()
			return
		}

		// Verificar acceso temporal
		var tempAccess models.TemporaryAccess
		if err := configs.DB.Where("user_id = ? AND entity_id = ? AND entity_type = ? AND expires_at > ?", userID, entityID, entityType, time.Now()).First(&tempAccess).Error; err == nil {
			c.Next()
			return
		}

		// Verificar si el área del usuario coincide con el área de la entidad
		if user.Area != entityArea {
			c.JSON(http.StatusForbidden, gin.H{
				"error":      "You do not have access to this resource",
				"entityArea": entityArea,
				"entityName": entityName,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
