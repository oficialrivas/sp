package controllers

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/utils"

	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"gorm.io/gorm"
)

// CreateDireccion crea una nueva dirección
// @Summary Crea una nueva dirección
// @Accept json
// @Produce json
// @Tags Direccion
// @Param Authorization header string true "Bearer token"
// @Param direccion body models.Direccion true "Datos de la dirección a crear"
// @Success 201 {object} models.Direccion
// @Router /direcciones [post]
func CreateDireccion(c *gin.Context) {
	var direccion models.Direccion
	if err := c.ShouldBindJSON(&direccion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	direccion.Area = claims.Area

	// Obtener el userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	direccion.UserID = uuid.MustParse(userID.(string))

	if err := configs.DB.Create(&direccion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la dirección"})
		return
	}

	c.JSON(http.StatusCreated, direccion)
}

// GetDireccionByID obtiene una dirección por su ID
// @Summary Obtiene una dirección por su ID
// @Produce json
// @Tags Direccion
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la dirección"
// @Success 200 {object} models.Direccion
// @Router /direcciones/{id} [get]
func GetDireccionByID(c *gin.Context) {
	id := c.Param("id")

	var direccion models.Direccion
	if err := configs.DB.Where("id = ?", id).Preload("Usuarios").Preload("Empleados").First(&direccion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dirección no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la dirección"})
		return
	}

	c.JSON(http.StatusOK, direccion)
}

// UpdateDireccion actualiza una dirección existente por su ID
// @Summary Actualiza una dirección existente por su ID
// @Accept json
// @Produce json
// @Tags Direccion
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la dirección a actualizar"
// @Param direccion body models.Direccion true "Datos de la dirección a actualizar"
// @Success 200 {object} models.Direccion
// @Router /direcciones/{id} [put]
func UpdateDireccion(c *gin.Context) {
	id := c.Param("id")

	var direccion models.Direccion
	if err := configs.DB.Where("id = ?", id).First(&direccion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dirección no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la dirección"})
		return
	}

	if err := c.ShouldBindJSON(&direccion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configs.DB.Save(&direccion)

	c.JSON(http.StatusOK, direccion)
}

// DeleteDireccion elimina una dirección por su ID
// @Summary Elimina una dirección por su ID
// @Produce json
// @Tags Direccion
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la dirección a eliminar"
// @Success 204 "Dirección eliminada exitosamente"
// @Router /direcciones/{id} [delete]
func DeleteDireccion(c *gin.Context) {
	id := c.Param("id")

	var direccion models.Direccion
	if err := configs.DB.Where("id = ?", id).First(&direccion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dirección no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la dirección"})
		return
	}

	configs.DB.Delete(&direccion)

	c.Status(http.StatusNoContent)
}
