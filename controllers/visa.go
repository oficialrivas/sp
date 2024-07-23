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

// CreateVisa crea una nueva visa
// @Summary Crea una nueva visa
// @Accept json
// @Produce json
// @Tags Visa
// @Param Authorization header string true "Bearer token"
// @Param valoracion body string true "Valoración de la visa"
// @Param pais body string true "País de la visa"
// @Param tipo body string true "Tipo de visa"
// @Param codigo body string true "Código único de la visa"
// @Param representante_id body string true "ID del representante"
// @Param nivel body string true "Nivel de la visa"
// @Param user_id body string true "ID del usuario"
// @Success 201 {object} models.Visa
// @Router /visas [post]
func CreateVisa(c *gin.Context) {
	var visa models.Visa
	if err := c.ShouldBindJSON(&visa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	visa.UserID = uuid.MustParse(userID.(string))

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
	visa.Area = claims.Area

	visa.ID = uuid.New()
	visa.CreatedAt = time.Now().UTC()
	visa.UpdatedAt = visa.CreatedAt

	if err := configs.DB.Create(&visa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la visa"})
		return
	}

	c.JSON(http.StatusCreated, visa)
}

// GetVisaByID obtiene una visa por su ID
// @Summary Obtiene una visa por su ID
// @Produce json
// @Tags Visa
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la visa"
// @Success 200 {object} models.Visa
// @Router /visas/{id} [get]
func GetVisaByID(c *gin.Context) {
	id := c.Param("id")

	var visa models.Visa
	if err := configs.DB.First(&visa, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Visa no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la visa"})
		return
	}

	c.JSON(http.StatusOK, visa)
}

// UpdateVisa actualiza una visa existente por su ID
// @Summary Actualiza una visa existente por su ID
// @Accept json
// @Produce json
// @Tags Visa
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la visa a actualizar"
// @Param valoracion body string true "Valoración de la visa"
// @Param pais body string true "País de la visa"
// @Param tipo body string true "Tipo de visa"
// @Param codigo body string true "Código único de la visa"
// @Param representante_id body string true "ID del representante"
// @Param nivel body string true "Nivel de la visa"
// @Param user_id body string true "ID del usuario"
// @Success 200 {object} models.Visa
// @Router /visas/{id} [put]
func UpdateVisa(c *gin.Context) {
	id := c.Param("id")

	var visa models.Visa
	if err := configs.DB.First(&visa, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Visa no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la visa"})
		return
	}

	if err := c.ShouldBindJSON(&visa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visa.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&visa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la visa"})
		return
	}

	c.JSON(http.StatusOK, visa)
}

// DeleteVisa elimina una visa por su ID
// @Summary Elimina una visa por su ID
// @Produce json
// @Tags Visa
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la visa a eliminar"
// @Success 204 "Visa eliminada exitosamente"
// @Router /visas/{id} [delete]
func DeleteVisa(c *gin.Context) {
	id := c.Param("id")

	var visa models.Visa
	if err := configs.DB.First(&visa, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Visa no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la visa"})
		return
	}

	if err := configs.DB.Delete(&visa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la visa"})
		return
	}

	c.Status(http.StatusNoContent)
}
