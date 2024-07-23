package controllers

import (
	"net/http"
	"time"
	"github.com/oficialrivas/sgi/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"gorm.io/gorm"
)

// CreateRedes crea una nueva entrada de red
// @Summary Crea una nueva entrada de red
// @Accept json
// @Produce json
// @Tags Redes
// @Param redes body models.Redes true "Datos de la red a crear"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} models.Redes
// @Router /redes [post]
func CreateRedes(c *gin.Context) {
	var redes models.Redes
	if err := c.ShouldBindJSON(&redes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	redes.UserID = uuid.MustParse(userID.(string))

	// Obtener el token desde el encabezado
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not found"})
		return
	}

	// Validar el token y obtener las claims
	claims, err := utils.ValidateJWT(tokenString, false)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
		return
	}

	// Asignar el área del usuario desde el token al correo
	redes.Area = claims.Area

	redes.ID = uuid.New()
	redes.CreatedAt = time.Now().UTC()
	redes.UpdatedAt = redes.CreatedAt

	if err := configs.DB.Create(&redes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la red"})
		return
	}

	c.JSON(http.StatusCreated, redes)
}

// GetRedesByID obtiene una entrada de red por su ID
// @Summary Obtiene una entrada de red por su ID
// @Produce json
// @Tags Redes
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la red"
// @Success 200 {object} models.Redes
// @Router /redes/{id} [get]
func GetRedesByID(c *gin.Context) {
	id := c.Param("id")

	var redes models.Redes
	if err := configs.DB.First(&redes, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Red no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la red"})
		return
	}

	c.JSON(http.StatusOK, redes)
}

// UpdateRedes actualiza una entrada de red existente por su ID
// @Summary Actualiza una entrada de red existente por su ID
// @Accept json
// @Produce json
// @Tags Redes
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la red a actualizar"
// @Param redes body models.Redes true "Datos de la red a actualizar"
// @Success 200 {object} models.Redes
// @Router /redes/{id} [put]
func UpdateRedes(c *gin.Context) {
	id := c.Param("id")

	var redes models.Redes
	if err := configs.DB.First(&redes, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Red no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la red"})
		return
	}

	if err := c.ShouldBindJSON(&redes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redes.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&redes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la red"})
		return
	}

	c.JSON(http.StatusOK, redes)
}

// DeleteRedes elimina una entrada de red por su ID
// @Summary Elimina una entrada de red por su ID
// @Produce json
// @Tags Redes
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la red a eliminar"
// @Success 204 "Red eliminada exitosamente"
// @Router /redes/{id} [delete]
func DeleteRedes(c *gin.Context) {
	id := c.Param("id")

	var redes models.Redes
	if err := configs.DB.First(&redes, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Red no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la red"})
		return
	}

	if err := configs.DB.Delete(&redes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la red"})
		return
	}

	c.Status(http.StatusNoContent)
}
