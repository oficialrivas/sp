package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
	"github.com/oficialrivas/sgi/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"gorm.io/gorm"
)

// CreatePasaporte crea un nuevo pasaporte
// @Summary Crea un nuevo pasaporte
// @Accept multipart/form-data
// @Produce json
// @Tags Pasaporte
// @Param Authorization header string true "Bearer token"
// @Param numero formData string true "Número del pasaporte"
// @Param foto formData file false "Foto del pasaporte"
// @Param pais formData string true "País del pasaporte"
// @Param tipo formData string true "Tipo del pasaporte"
// @Param codigo formData string true "Código único del pasaporte"
// @Param representante_id formData string true "ID del representante"
// @Param nivel formData string true "Nivel del pasaporte"
// @Param user_id formData string true "ID del usuario"
// @Success 201 {object} models.Pasaporte
// @Router /pasaportes [post]
func CreatePasaporte(c *gin.Context) {
	var pasaporte models.Pasaporte
	pasaporte.Numero = c.PostForm("numero")
	pasaporte.Pais = c.PostForm("pais")
	pasaporte.Tipo = c.PostForm("tipo")
	pasaporte.Codigo = c.PostForm("codigo")
	pasaporte.RepresentanteID = uuid.MustParse(c.PostForm("representante_id"))
	pasaporte.UserID = uuid.MustParse(c.PostForm("user_id"))

	// Guardar la foto si se envía
	file, err := c.FormFile("foto")
	if err == nil {
		filename := uuid.New().String() + filepath.Ext(file.Filename)
		filepath := filepath.Join("static", filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	pasaporte.UserID = uuid.MustParse(userID.(string))

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
	pasaporte.Area = claims.Area

		pasaporte.Foto = filename
	}

	pasaporte.ID = uuid.New()
	pasaporte.CreatedAt = time.Now().UTC()
	pasaporte.UpdatedAt = pasaporte.CreatedAt

	if err := configs.DB.Create(&pasaporte).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el pasaporte"})
		return
	}

	c.JSON(http.StatusCreated, pasaporte)
}

// GetPasaporteByID obtiene un pasaporte por su ID
// @Param Authorization header string true "Bearer token"
// @Summary Obtiene un pasaporte por su ID
// @Produce json
// @Tags Pasaporte
// @Param id path string true "ID del pasaporte"
// @Success 200 {object} models.Pasaporte
// @Router /pasaportes/{id} [get]
func GetPasaporteByID(c *gin.Context) {
	id := c.Param("id")

	var pasaporte models.Pasaporte
	if err := configs.DB.First(&pasaporte, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pasaporte no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el pasaporte"})
		return
	}

	c.JSON(http.StatusOK, pasaporte)
}

// UpdatePasaporte actualiza un pasaporte existente por su ID
// @Summary Actualiza un pasaporte existente por su ID
// @Accept multipart/form-data
// @Produce json
// @Tags Pasaporte
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del pasaporte a actualizar"
// @Param numero formData string true "Número del pasaporte"
// @Param foto formData file false "Foto del pasaporte"
// @Param pais formData string true "País del pasaporte"
// @Param tipo formData string true "Tipo del pasaporte"
// @Param codigo formData string true "Código único del pasaporte"
// @Param representante_id formData string true "ID del representante"
// @Param nivel formData string true "Nivel del pasaporte"
// @Param user_id formData string true "ID del usuario"
// @Success 200 {object} models.Pasaporte
// @Router /pasaportes/{id} [put]
func UpdatePasaporte(c *gin.Context) {
	id := c.Param("id")

	var pasaporte models.Pasaporte
	if err := configs.DB.First(&pasaporte, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pasaporte no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el pasaporte"})
		return
	}

	pasaporte.Numero = c.PostForm("numero")
	pasaporte.Pais = c.PostForm("pais")
	pasaporte.Tipo = c.PostForm("tipo")
	pasaporte.Codigo = c.PostForm("codigo")
	pasaporte.RepresentanteID = uuid.MustParse(c.PostForm("representante_id"))
	pasaporte.UserID = uuid.MustParse(c.PostForm("user_id"))

	// Guardar la nueva foto si se envía
	file, err := c.FormFile("foto")
	if err == nil {
		// Eliminar la foto anterior
		if pasaporte.Foto != "" {
			oldFilepath := filepath.Join("static", pasaporte.Foto)
			os.Remove(oldFilepath)
		}

		filename := uuid.New().String() + filepath.Ext(file.Filename)
		filepath := filepath.Join("static", filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		pasaporte.Foto = filename
	}

	pasaporte.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&pasaporte).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el pasaporte"})
		return
	}

	c.JSON(http.StatusOK, pasaporte)
}

// DeletePasaporte elimina un pasaporte por su ID
// @Summary Elimina un pasaporte por su ID
// @Produce json
// @Tags Pasaporte
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del pasaporte a eliminar"
// @Success 204 "Pasaporte eliminado exitosamente"
// @Router /pasaportes/{id} [delete]
func DeletePasaporte(c *gin.Context) {
	id := c.Param("id")

	var pasaporte models.Pasaporte
	if err := configs.DB.First(&pasaporte, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pasaporte no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el pasaporte"})
		return
	}

	// Eliminar la foto asociada si existe
	if pasaporte.Foto != "" {
		filepath := filepath.Join("static", pasaporte.Foto)
		os.Remove(filepath)
	}

	if err := configs.DB.Delete(&pasaporte).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el pasaporte"})
		return
	}

	c.Status(http.StatusNoContent)
}
