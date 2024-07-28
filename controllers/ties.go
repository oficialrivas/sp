package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"

	"gorm.io/gorm"
)

// CreateTie crea una nueva TIE
// @Summary Crea una nueva TIE
// @Accept json
// @Produce json
// @Tags Tie
// @Param Authorization header string true "Bearer token"
// @Param TIE body models.Tie true "Datos de la TIE a crear"
// @Success 201 {object} models.Tie
// @Router /ties [post]
func CreateTie(c *gin.Context) {
	var input struct {
		Nombre    string        `json:"nombre" binding:"required"`
		Modalidad []struct {
			ID uuid.UUID `json:"id"`
		} `json:"modalidad" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	tie := models.Tie{
		Nombre: input.Nombre,
		UserID: uuid.MustParse(userID.(string)),
	}

	for _, modalidadID := range input.Modalidad {
		var modalidad models.Modalidad
		if err := configs.DB.First(&modalidad, "id = ?", modalidadID.ID).Error; err == nil {
			tie.Modalidad = append(tie.Modalidad, modalidad)
		}
	}

	tie.ID = uuid.New()
	tie.CreatedAt = time.Now().UTC()
	tie.UpdatedAt = tie.CreatedAt

	if err := configs.DB.Create(&tie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la TIE"})
		return
	}

	c.JSON(http.StatusCreated, tie)
}

// GetTieByID obtiene una TIE por su ID
// @Summary Obtiene una TIE por su ID
// @Description Obtiene una TIE por su ID
// @Tags Tie
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la TIE"
// @Success 200 {object} models.Tie
// @Failure 404 {object} models.ErrorResponse "TIE no encontrada"
// @Failure 500 {object} models.ErrorResponse "Error al buscar la TIE"
// @Router /ties/{id} [get]
func GetTieByID(c *gin.Context) {
	id := c.Param("id")

	var tie models.Tie
	if err := configs.DB.Preload("Modalidad").First(&tie, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "TIE no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la TIE"})
		return
	}

	c.JSON(http.StatusOK, tie)
}

// UpdateTie actualiza una TIE existente por su ID
// @Summary Actualiza una TIE existente por su ID
// @Accept json
// @Produce json
// @Tags Tie
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la TIE a actualizar"
// @Param nombre body string true "Nombre de la TIE"
// @Param area body string true "Área de la TIE"
// @Param user_id body string true "ID del usuario"
// @Param vehiculo body models.Tie true "Datos de la TIE a actualizar"
// @Success 200 {object} models.Tie
// @Router /ties/{id} [put]
func UpdateTie(c *gin.Context) {
	id := c.Param("id")

	var tie models.Tie
	if err := configs.DB.First(&tie, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "TIE no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la TIE"})
		return
	}

	if err := c.ShouldBindJSON(&tie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tie.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&tie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la TIE"})
		return
	}

	c.JSON(http.StatusOK, tie)
}

// DeleteTie elimina una TIE por su ID
// @Summary Elimina una TIE por su ID
// @Produce json
// @Tags Tie
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la TIE a eliminar"
// @Success 204 "TIE eliminada exitosamente"
// @Router /ties/{id} [delete]
func DeleteTie(c *gin.Context) {
	id := c.Param("id")

	var tie models.Tie
	if err := configs.DB.First(&tie, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "TIE no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la TIE"})
		return
	}

	if err := configs.DB.Delete(&tie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la TIE"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTies obtiene todas las TIEs
// @Summary Obtiene todas las TIEs
// @Produce json
// @Tags Tie
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Tie
// @Router /ties [get]
func GetTies(c *gin.Context) {
	var ties []models.Tie
	if err := configs.DB.Preload("Modalidad").Find(&ties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las TIEs"})
		return
	}

	c.JSON(http.StatusOK, ties)
}