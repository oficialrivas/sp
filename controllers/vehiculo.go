package controllers

import (
	"net/http"
	"github.com/oficialrivas/sgi/utils"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/config"
	"gorm.io/gorm"
)

// CreateVehiculo crea un nuevo vehículo
// @Summary Crea un nuevo vehículo
// @Accept json
// @Produce json
// @Tags Vehiculo
// @Param Authorization header string true "Bearer token"
// @Param vehiculo body models.Vehiculo true "Datos del vehículo a crear"
// @Success 201 {object} models.Vehiculo
// @Router /vehiculos [post]
func CreateVehiculo(c *gin.Context) {
	var vehiculo models.Vehiculo
	if err := c.ShouldBindJSON(&vehiculo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	vehiculo.UserID = uuid.MustParse(userID.(string))

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
	vehiculo.Area = claims.Area

	if err := configs.DB.Create(&vehiculo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el vehículo"})
		return
	}

	c.JSON(http.StatusCreated, vehiculo)
}

// GetVehiculoByID obtiene un vehículo por su ID
// @Summary Obtiene un vehículo por su ID
// @Produce json
// @Tags Vehiculo
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del vehículo"
// @Success 200 {object} models.Vehiculo
// @Router /vehiculos/{id} [get]
func GetVehiculoByID(c *gin.Context) {
	id := c.Param("id")

	var vehiculo models.Vehiculo
	if err := configs.DB.Where("id = ?", id).Preload("Usuarios").First(&vehiculo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el vehículo"})
		return
	}

	c.JSON(http.StatusOK, vehiculo)
}

// GetVehiculoByMatricula obtiene un vehículo por su matrícula
// @Summary Obtiene un vehículo por su matrícula
// @Produce json
// @Tags Vehiculo
// @Param Authorization header string true "Bearer token"
// @Param matricula query string true "Matrícula del vehículo"
// @Success 200 {object} models.Vehiculo
// @Router /vehiculos/search [get]
func GetVehiculoByMatricula(c *gin.Context) {
	matricula := c.Query("matricula")

	var vehiculo models.Vehiculo
	if err := configs.DB.Where("matricula = ?", matricula).Preload("Usuarios").First(&vehiculo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el vehículo"})
		return
	}

	c.JSON(http.StatusOK, vehiculo)
}

// UpdateVehiculo actualiza un vehículo existente por su ID
// @Summary Actualiza un vehículo existente por su ID
// @Accept json
// @Produce json
// @Tags Vehiculo
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del vehículo a actualizar"
// @Param vehiculo body models.Vehiculo true "Datos del vehículo a actualizar"
// @Success 200 {object} models.Vehiculo
// @Router /vehiculos/{id} [put]
func UpdateVehiculo(c *gin.Context) {
	id := c.Param("id")

	var vehiculo models.Vehiculo
	if err := configs.DB.Where("id = ?", id).First(&vehiculo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el vehículo"})
		return
	}

	if err := c.ShouldBindJSON(&vehiculo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configs.DB.Save(&vehiculo)

	c.JSON(http.StatusOK, vehiculo)
}

// DeleteVehiculo elimina un vehículo por su ID
// @Summary Elimina un vehículo por su ID
// @Produce json
// @Tags Vehiculo
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del vehículo a eliminar"
// @Success 204 "Vehículo eliminado exitosamente"
// @Router /vehiculos/{id} [delete]
func DeleteVehiculo(c *gin.Context) {
	id := c.Param("id")

	var vehiculo models.Vehiculo
	if err := configs.DB.Where("id = ?", id).First(&vehiculo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el vehículo"})
		return
	}

	configs.DB.Delete(&vehiculo)

	c.Status(http.StatusNoContent)
}
