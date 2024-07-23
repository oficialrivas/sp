package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/utils"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"gorm.io/gorm"
)

// CreateDocumento crea un nuevo documento
// @Summary Crea un nuevo documento
// @Accept multipart/form-data
// @Produce json
// @Tags Documento
// @Param Authorization header string true "Bearer token"
// @Param numero formData string true "Número del documento"
// @Param documento formData file false "Archivo del documento"
// @Param nombre formData string true "Nombre del documento"
// @Param tipo formData string true "Tipo del documento"
// @Param codigo formData string true "Código único del documento"
// @Param nivel formData string true "Nivel del documento"
// @Param user_id formData string true "ID del usuario"
// @Success 201 {object} models.Documento
// @Router /documentos [post]
func CreateDocumento(c *gin.Context) {
	var documento models.Documento
	documento.Numero = c.PostForm("numero")
	documento.Nombre = c.PostForm("nombre")
	documento.Tipo = c.PostForm("tipo")
	documento.Codigo = c.PostForm("codigo")
	documento.UserID = uuid.MustParse(c.PostForm("user_id"))

	// Guardar el archivo del documento si se envía
	file, err := c.FormFile("documento")
	if err == nil {
		filename := uuid.New().String() + filepath.Ext(file.Filename)
		filepath := filepath.Join("static/documentos", filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
			return
		}
	
		documento.UserID = uuid.MustParse(userID.(string))

		documento.Documento = filename
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
	documento.Area = claims.Area

	documento.ID = uuid.New()
	documento.CreatedAt = time.Now().UTC()
	documento.UpdatedAt = documento.CreatedAt

	if err := configs.DB.Create(&documento).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el documento"})
		return
	}

	c.JSON(http.StatusCreated, documento)
}

// GetDocumentoByID obtiene un documento por su ID
// @Summary Obtiene un documento por su ID
// @Produce json
// @Tags Documento
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del documento"
// @Success 200 {object} models.Documento
// @Router /documentos/{id} [get]
func GetDocumentoByID(c *gin.Context) {
	id := c.Param("id")

	var documento models.Documento
	if err := configs.DB.First(&documento, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el documento"})
		return
	}

	c.JSON(http.StatusOK, documento)
}

// UpdateDocumento actualiza un documento existente por su ID
// @Summary Actualiza un documento existente por su ID
// @Accept multipart/form-data
// @Produce json
// @Tags Documento
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del documento a actualizar"
// @Param numero formData string true "Número del documento"
// @Param documento formData file false "Archivo del documento"
// @Param nombre formData string true "Nombre del documento"
// @Param tipo formData string true "Tipo del documento"
// @Param codigo formData string true "Código único del documento"
// @Param nivel formData string true "Nivel del documento"
// @Param user_id formData string true "ID del usuario"
// @Success 200 {object} models.Documento
// @Router /documentos/{id} [put]
func UpdateDocumento(c *gin.Context) {
	id := c.Param("id")

	var documento models.Documento
	if err := configs.DB.First(&documento, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el documento"})
		return
	}

	documento.Numero = c.PostForm("numero")
	documento.Nombre = c.PostForm("nombre")
	documento.Tipo = c.PostForm("tipo")
	documento.Codigo = c.PostForm("codigo")
	documento.UserID = uuid.MustParse(c.PostForm("user_id"))

	// Guardar el nuevo archivo del documento si se envía
	file, err := c.FormFile("documento")
	if err == nil {
		// Eliminar el archivo anterior
		if documento.Documento != "" {
			oldFilepath := filepath.Join("static/documentos", documento.Documento)
			os.Remove(oldFilepath)
		}

		filename := uuid.New().String() + filepath.Ext(file.Filename)
		filepath := filepath.Join("static/documentos", filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		documento.Documento = filename
	}

	documento.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&documento).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el documento"})
		return
	}

	c.JSON(http.StatusOK, documento)
}

// DeleteDocumento elimina un documento por su ID
// @Summary Elimina un documento por su ID
// @Produce json
// @Tags Documento
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del documento a eliminar"
// @Success 204 "Documento eliminado exitosamente"
// @Router /documentos/{id} [delete]
func DeleteDocumento(c *gin.Context) {
	id := c.Param("id")

	var documento models.Documento
	if err := configs.DB.First(&documento, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el documento"})
		return
	}

	// Eliminar el archivo asociado si existe
	if documento.Documento != "" {
		filepath := filepath.Join("static/documentos", documento.Documento)
		os.Remove(filepath)
	}

	if err := configs.DB.Delete(&documento).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el documento"})
		return
	}

	c.Status(http.StatusNoContent)
}
