package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/utils"
	"gorm.io/gorm"
)

// CreateModalidad crea una nueva modalidad
// @Summary Crea una nueva modalidad
// @Accept json
// @Produce json
// @Tags Modalidad
// @Param Authorization header string true "Bearer token"
// @Param modalidad body models.Modalidad true "Datos de la Modalidad a crear"
// @Success 201 {object} models.Modalidad
// @Router /modalidades [post]
func CreateModalidad(c *gin.Context) {
	var modalidad models.Modalidad
	if err := c.ShouldBindJSON(&modalidad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	modalidad.UserID = uuid.MustParse(userID.(string))

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
	modalidad.Area = claims.Area

	modalidad.ID = uuid.New()
	modalidad.CreatedAt = time.Now().UTC()
	modalidad.UpdatedAt = modalidad.CreatedAt

	if err := configs.DB.Create(&modalidad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la modalidad"})
		return
	}

	c.JSON(http.StatusCreated, modalidad)
}

// GetModalidadByID obtiene una modalidad por su ID
// @Summary Obtiene una modalidad por su ID
// @Produce json
// @Tags Modalidad
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la modalidad"
// @Success 200 {object} models.Modalidad
// @Router /modalidades/{id} [get]
func GetModalidadByID(c *gin.Context) {
	id := c.Param("id")

	var modalidad models.Modalidad
	if err := configs.DB.First(&modalidad, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Modalidad no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la modalidad"})
		return
	}

	c.JSON(http.StatusOK, modalidad)
}

// UpdateModalidad actualiza una modalidad existente por su ID
// @Summary Actualiza una modalidad existente por su ID
// @Accept json
// @Produce json
// @Tags Modalidad
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la modalidad a actualizar"
// @Param nombre body string true "Nombre de la modalidad"
// @Param user_id body string true "ID del usuario"
// @Param vehiculo body models.Modalidad true "Datos de la Modalidad a actualizar"
// @Success 200 {object} models.Modalidad
// @Router /modalidades/{id} [put]
func UpdateModalidad(c *gin.Context) {
	id := c.Param("id")

	var modalidad models.Modalidad
	if err := configs.DB.First(&modalidad, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Modalidad no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la modalidad"})
		return
	}

	if err := c.ShouldBindJSON(&modalidad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modalidad.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&modalidad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la modalidad"})
		return
	}

	c.JSON(http.StatusOK, modalidad)
}

// DeleteModalidad elimina una modalidad por su ID
// @Summary Elimina una modalidad por su ID
// @Produce json
// @Tags Modalidad
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la modalidad a eliminar"
// @Success 204 "Modalidad eliminada exitosamente"
// @Router /modalidades/{id} [delete]
func DeleteModalidad(c *gin.Context) {
	id := c.Param("id")

	var modalidad models.Modalidad
	if err := configs.DB.First(&modalidad, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Modalidad no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la modalidad"})
		return
	}

	if err := configs.DB.Delete(&modalidad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la modalidad"})
		return
	}

	c.Status(http.StatusNoContent)
}


// GetAllModalidades obtiene todas las modalidades
// @Summary Obtiene todas las modalidades
// @Produce json
// @Tags Modalidad
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Modalidad
// @Router /modalidades [get]
func GetAllModalidades(c *gin.Context) {
	var modalidades []models.Modalidad
	if err := configs.DB.Find(&modalidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las modalidades"})
		return
	}
	c.JSON(http.StatusOK, modalidades)
}