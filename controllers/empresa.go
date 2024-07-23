package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/utils"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"gorm.io/gorm"
)

// CreateEmpresa crea una nueva empresa
// @Summary Crea una nueva empresa
// @Accept json
// @Produce json
// @Tags Empresa
// @Param Authorization header string true "Bearer token"
// @Param empresa body models.Empresa true "Datos de la empresa a crear"
// @Success 201 {object} models.Empresa
// @Router /empresas [post]
// @Security BearerAuth
func CreateEmpresa(c *gin.Context) {
	var empresa models.Empresa
	if err := c.ShouldBindJSON(&empresa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	empresa.UserID = uuid.MustParse(userID.(string))

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
		empresa.Area = claims.Area

	if err := configs.DB.Create(&empresa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la empresa"})
		return
	}

	c.JSON(http.StatusCreated, empresa)
}

// GetEmpresaByID obtiene una empresa por su ID
// @Summary Obtiene una empresa por su ID
// @Produce json
// @Tags Empresa
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la empresa"
// @Success 200 {object} models.Empresa
// @Router /empresas/{id} [get]
func GetEmpresaByID(c *gin.Context) {
	id := c.Param("id")

	var empresa models.Empresa
	if err := configs.DB.Where("id = ?", id).Preload("Socios").Preload("Empleados").First(&empresa).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la empresa"})
		return
	}

	c.JSON(http.StatusOK, empresa)
}

// GetEmpresaByRIF obtiene una empresa por su RIF
// @Summary Obtiene una empresa por su RIF
// @Produce json
// @Tags Empresa
// @Param Authorization header string true "Bearer token"
// @Param rif query string true "RIF de la empresa"
// @Success 200 {object} models.Empresa
// @Router /empresas/search [get]
func GetEmpresaByRIF(c *gin.Context) {
	rif := c.Query("rif")

	var empresa models.Empresa
	if err := configs.DB.Where("rif = ?", rif).Preload("Socios").Preload("Empleados").First(&empresa).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la empresa"})
		return
	}

	c.JSON(http.StatusOK, empresa)
}

// UpdateEmpresa actualiza una empresa existente por su ID
// @Summary Actualiza una empresa existente por su ID
// @Accept json
// @Produce json
// @Tags Empresa
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la empresa a actualizar"
// @Param empresa body models.Empresa true "Datos de la empresa a actualizar"
// @Success 200 {object} models.Empresa
// @Router /empresas/{id} [put]
func UpdateEmpresa(c *gin.Context) {
	id := c.Param("id")

	var empresa models.Empresa
	if err := configs.DB.Where("id = ?", id).First(&empresa).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la empresa"})
		return
	}

	if err := c.ShouldBindJSON(&empresa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configs.DB.Save(&empresa)

	c.JSON(http.StatusOK, empresa)
}

// DeleteEmpresa elimina una empresa por su ID
// @Summary Elimina una empresa por su ID
// @Produce json
// @Tags Empresa
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la empresa a eliminar"
// @Success 204 "Empresa eliminada exitosamente"
// @Router /empresas/{id} [delete]
func DeleteEmpresa(c *gin.Context) {
	id := c.Param("id")

	var empresa models.Empresa
	if err := configs.DB.Where("id = ?", id).First(&empresa).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la empresa"})
		return
	}

	configs.DB.Delete(&empresa)

	c.Status(http.StatusNoContent)
}
