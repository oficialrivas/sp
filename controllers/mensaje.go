package controllers

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/utils"
)

// CreateMensaje crea un nuevo registro de Mensaje
// @Summary Crea un nuevo registro de Mensaje
// @Description Crea un nuevo registro de Mensaje con los datos proporcionados
// @Tags mensaje
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param descripcion formData string true "Descripción del Mensaje"
// @Param fecha formData string true "Fecha del Mensaje (YYYY-MM-DD)"
// @Param lugar formData string true "Lugar del Mensaje"
// @Param modalidad formData string true "Modalidad del Mensaje"
// @Param nombre formData string true "Nombre del Mensaje"
// @Param parroquia formData string true "Parroquia del Mensaje"
// @Param redi formData string true "REDI del Mensaje"
// @Param zodi formData string true "ZODI del Mensaje"
// @Param adi formData string true "ADI del Mensaje"
// @Param tie formData string true "TIE del Mensaje"
// @Param area formData string true "Área del Mensaje"
// @Param procesado formData bool true "Procesado"
// @Param imagen formData file false "Imagen del Mensaje"
// @Success 200 {object} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /mensajes [post]
// @Security ApiKeyAuth
func CreateMensaje(c *gin.Context) {
	var input struct {
		Descripcion  string `form:"descripcion" binding:"required"`
		Fecha        string `form:"fecha" binding:"required"`
		Lugar        string `form:"lugar" binding:"required"`
		Modalidad    string `form:"modalidad" binding:"required"`
		Nombre       string `form:"nombre" binding:"required"`
		Parroquia    string `form:"parroquia" binding:"required"`
		REDI         string `form:"redi" binding:"required"`
		ZODI         string `form:"zodi" binding:"required"`
		ADI          string `form:"adi" binding:"required"`
		Tie          string `form:"tie" binding:"required"`
		Area         string `form:"area" binding:"required"`
		Procesado    bool   `form:"procesado" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.Fecha)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
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

	// Manejar la carga de la imagen
	file, err := c.FormFile("imagen")
	var imagenURL string
	if err == nil {
		// Guardar la imagen en una carpeta específica
		imagenPath := filepath.Join("static", file.Filename)
		if err := c.SaveUploadedFile(file, imagenPath); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to upload image"})
			return
		}
		imagenURL = imagenPath
	}

	mensaje := models.Mensaje{
		Descripcion:  input.Descripcion,
		Fecha:        parsedDate,
		Lugar:        input.Lugar,
		Modalidad:    input.Modalidad,
		Nombre:       input.Nombre,
		Parroquia:    input.Parroquia,
		REDI:         input.REDI,
		ZODI:         input.ZODI,
		ADI:          input.ADI,
		Tie:          input.Tie,
		Area:         claims.Area, // Asignar el área del usuario desde el token
		UserID:       uuid.MustParse(userID.(string)),
		ImagenURL:    imagenURL,
		Procesado:    input.Procesado,
	
	}

	if err := configs.DB.Create(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mensaje)
}

// GetMensaje obtiene un Mensaje por su ID
// @Summary Obtiene un Mensaje por su ID
// @Description Obtiene los datos de un Mensaje por su ID
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del Mensaje"
// @Success 200 {object} models.Mensaje
// @Failure 404 {object} models.ErrorResponse
// @Router /mensajes/{id} [get]
// @Security ApiKeyAuth
func GetMensaje(c *gin.Context) {
	id := c.Param("id")
	var mensaje models.Mensaje
	if err := configs.DB.Where("id = ?", id).First(&mensaje).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Mensaje no encontrado"})
		return
	}

	c.JSON(http.StatusOK, mensaje)
}

// UpdateMensaje actualiza un Mensaje existente
// @Summary Actualiza un Mensaje existente
// @Description Actualiza los datos de un Mensaje existente
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del Mensaje"
// @Param mensaje body models.Mensaje true "Datos del Mensaje actualizados"
// @Success 200 {object} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /mensajes/{id} [put]
// @Security ApiKeyAuth
func UpdateMensaje(c *gin.Context) {
	id := c.Param("id")
	var mensaje models.Mensaje
	if err := configs.DB.Where("id = ?", id).First(&mensaje).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Mensaje no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&mensaje); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	mensaje.ID, _ = uuid.Parse(id)
	mensaje.UpdatedAt = time.Now().UTC()
	if err := configs.DB.Save(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mensaje)
}

// DeleteMensaje elimina un Mensaje por su ID
// @Summary Elimina un Mensaje por su ID
// @Description Elimina un Mensaje por su ID
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del Mensaje"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /mensajes/{id} [delete]
// @Security ApiKeyAuth
func DeleteMensaje(c *gin.Context) {
	id := c.Param("id")
	var mensaje models.Mensaje
	if err := configs.DB.Where("id = ?", id).First(&mensaje).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Mensaje no encontrado"})
		return
	}

	if err := configs.DB.Delete(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetMensajes obtiene todos los Mensajes
// @Summary Obtiene todos los Mensajes
// @Description Obtiene todos los registros de Mensaje
// @Tags mensaje
// @Accept json
// @Produce json
// @Success 200 {array} models.Mensaje
// @Failure 500 {object} models.ErrorResponse
// @Router /mensajes [get]
// @Security ApiKeyAuth
func GetMensajes(c *gin.Context) {
	var mensajes []models.Mensaje
	if err := configs.DB.Find(&mensajes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mensajes)
}

// FilterMensajes obtiene Mensajes basados en filtros opcionales
// @Summary Obtiene Mensajes filtrados
// @Description Obtiene registros de Mensaje en un periodo específico, y filtrados por tie, redi, zodi y modalidad
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Param tie query string false "TIE"
// @Param redi query string false "REDI"
// @Param zodi query string false "ZODI"
// @Param adi query string false "ADI"
// @Param modalidad query string false "Modalidad"
// @Success 200 {array} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /mensajes/filter [get]
// @Security ApiKeyAuth
func FilterMensajes(c *gin.Context) {
	var mensajes []models.Mensaje

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tie := c.Query("tie")
	redi := c.Query("redi")
	zodi := c.Query("zodi")
	adi := c.Query("adi")
	modalidad := c.Query("modalidad")

	db := configs.DB

	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid start_date format. Use YYYY-MM-DD."})
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid end_date format. Use YYYY-MM-DD."})
			return
		}
		db = db.Where("fecha BETWEEN ? AND ?", start, end)
	}

	if tie != "" {
		db = db.Where("tie = ?", tie)
	}
	if redi != "" {
		db = db.Where("redi = ?", redi)
	}
	if zodi != "" {
		db = db.Where("zodi = ?", zodi)
	}
	if adi != "" {
		db = db.Where("adi = ?", adi)
	}
	if modalidad != "" {
		db = db.Where("modalidad = ?", modalidad)
	}

	if err := db.Find(&mensajes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mensajes)
}


// UpdateMensajeStatus actualiza el campo Procesado de un Mensaje a true
// @Summary Actualiza el campo Procesado de un Mensaje a true
// @Description Actualiza el campo Procesado de un Mensaje a true
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del Mensaje"
// @Success 200 {object} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /mensajes/{id}/procesado [put]
// @Security ApiKeyAuth
func UpdateMensajeStatus(c *gin.Context) {
	id := c.Param("id")

	// Verificar si el ID es válido
	mensajeID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	var mensaje models.Mensaje
	if err := configs.DB.Where("id = ?", mensajeID).First(&mensaje).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Mensaje no encontrado"})
		return
	}

	mensaje.Procesado = true

	if err := configs.DB.Save(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mensaje)
}