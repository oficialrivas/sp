package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/utils"
	"github.com/oficialrivas/sgi/models"
	"gorm.io/gorm"
)

// CreateCorreo crea un nuevo correo
// @Summary Crea un nuevo correo
// @Accept json
// @Produce json
// @Tags Correo
// @Param correo body models.Correo true "Datos de la dirección a crear"
// @Param Authorization header string true "Bearer token"
// @Param tipo body string true "Tipo del correo"
// @Param area body string true "Área del correo"
// @Param direccion body string true "Dirección del correo"
// @Param representante_id body string true "ID del representante"
// @Param user_id body string true "ID del usuario"
// @Success 201 {object} models.Correo
// @Router /correos [post]
func CreateCorreo(c *gin.Context) {
	var correo models.Correo
	if err := c.ShouldBindJSON(&correo); err != nil {
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
	correo.Area = claims.Area

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	correo.UserID = uuid.MustParse(userID.(string))

	correo.ID = uuid.New()
	correo.CreatedAt = time.Now().UTC()
	correo.UpdatedAt = correo.CreatedAt

	if err := configs.DB.Create(&correo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el correo"})
		return
	}

	c.JSON(http.StatusCreated, correo)
}

// GetCorreoByID obtiene un correo por su ID
// @Summary Obtiene un correo por su ID
// @Produce json
// @Tags Correo
// @Success 201 {object} models.Correo
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del correo"
// @Success 200 {object} models.Correo
// @Router /correos/{id} [get]
func GetCorreoByID(c *gin.Context) {
	id := c.Param("id")

	var correo models.Correo
	if err := configs.DB.First(&correo, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Correo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el correo"})
		return
	}

	c.JSON(http.StatusOK, correo)
}

// UpdateCorreo actualiza un correo existente por su ID
// @Summary Actualiza un correo existente por su ID
// @Accept json
// @Produce json
// @Tags Correo
// @Param direccion body models.Correo true "Datos del correo a actualizar"
// @Success 201 {object} models.Correo
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del correo a actualizar"
// @Param tipo body string true "Tipo del correo"
// @Param area body string true "Área del correo"
// @Param direccion body string true "Dirección del correo"
// @Param representante_id body string true "ID del representante"
// @Param user_id body string true "ID del usuario"
// @Success 200 {object} models.Correo
// @Router /correos/{id} [put]
func UpdateCorreo(c *gin.Context) {
	id := c.Param("id")

	var correo models.Correo
	if err := configs.DB.First(&correo, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Correo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el correo"})
		return
	}

	if err := c.ShouldBindJSON(&correo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	correo.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&correo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el correo"})
		return
	}

	c.JSON(http.StatusOK, correo)
}

// DeleteCorreo elimina un correo por su ID
// @Summary Elimina un correo por su ID
// @Produce json
// @Tags Correo
// @Success 201 {object} models.Correo
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del correo a eliminar"
// @Success 204 "Correo eliminado exitosamente"
// @Router /correos/{id} [delete]
func DeleteCorreo(c *gin.Context) {
	id := c.Param("id")

	var correo models.Correo
	if err := configs.DB.First(&correo, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Correo no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el correo"})
		return
	}

	if err := configs.DB.Delete(&correo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el correo"})
		return
	}

	c.Status(http.StatusNoContent)
}
