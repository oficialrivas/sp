package controllers

import (
	"net/http"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
)

type AccessRequest struct {
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	EntityID   uuid.UUID `json:"entity_id" binding:"required"`
	EntityType string    `json:"entity_type" binding:"required"`
	ExpiresAt  time.Time `json:"expires_at" binding:"required"`
}

// GrantTemporaryAccess otorga acceso temporal a un usuario
// @Summary Otorga acceso temporal a un usuario a una entidad específica
// @Description Otorga acceso temporal a un usuario de otra área a una entidad específica
// @Tags configuracion
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param access body AccessRequest true "Datos de acceso temporal"
// @Success 200 {object} models.TemporaryAccess
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /configuracion/acceso-temporal [post]
// @Security ApiKeyAuth
func GrantTemporaryAccess(c *gin.Context) {
	var request AccessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	access := models.TemporaryAccess{
		ID:         uuid.New(),
		UserID:     request.UserID,
		EntityID:   request.EntityID,
		EntityType: request.EntityType,
		ExpiresAt:  request.ExpiresAt,
	}

	if err := configs.DB.Create(&access).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, access)
}

// AddAreaRequest representa la solicitud para agregar un área
type AddAreaRequest struct {
	Area string `json:"area" binding:"required"`
}

// AddArea agrega un área nueva a la lista de áreas válidas y actualiza usuarios
// @Summary Agrega un área nueva
// @Description Agrega un área nueva a la lista de áreas válidas
// @Tags configuracion
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param area body AddAreaRequest true "Nombre del área"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /configuracion/area [post]
// @Security ApiKeyAuth
func AddArea(c *gin.Context) {
	var request AddAreaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Verificar si el área ya existe
	for _, area := range models.ValidAreas {
		if strings.EqualFold(area, request.Area) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "El área ya existe"})
			return
		}
	}

	// Agregar el área
	models.ValidAreas = append(models.ValidAreas, request.Area)
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Área agregada correctamente"})
}

// RemoveAreaRequest representa la solicitud para eliminar un área
type RemoveAreaRequest struct {
	Area string `json:"area" binding:"required"`
}

// RemoveArea elimina un área de la lista de áreas válidas
// @Summary Elimina un área
// @Description Elimina un área de la lista de áreas válidas
// @Tags configuracion
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param area body RemoveAreaRequest true "Nombre del área"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /configuracion/area [delete]
// @Security ApiKeyAuth
func RemoveArea(c *gin.Context) {
	var request RemoveAreaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Buscar y eliminar el área
	for i, area := range models.ValidAreas {
		if strings.EqualFold(area, request.Area) {
			models.ValidAreas = append(models.ValidAreas[:i], models.ValidAreas[i+1:]...)
			c.JSON(http.StatusOK, models.SuccessResponse{Message: "Área eliminada correctamente"})
			return
		}
	}

	c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "El área no existe"})
}

// UpdateAreaRequest representa la solicitud para actualizar un área
type UpdateAreaRequest struct {
	OldArea string `json:"old_area" binding:"required"`
	NewArea string `json:"new_area" binding:"required"`
}

// UpdateArea actualiza el nombre de un área en la lista de áreas válidas
// @Summary Actualiza un área
// @Description Actualiza el nombre de un área en la lista de áreas válidas
// @Tags configuracion
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param area body UpdateAreaRequest true "Datos del área"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /configuracion/area [put]
// @Security ApiKeyAuth
func UpdateArea(c *gin.Context) {
	var request UpdateAreaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Buscar y actualizar el área
	for i, area := range models.ValidAreas {
		if strings.EqualFold(area, request.OldArea) {
			models.ValidAreas[i] = request.NewArea
			c.JSON(http.StatusOK, models.SuccessResponse{Message: "Área actualizada correctamente"})
			return
		}
	}

	c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "El área no existe"})
}