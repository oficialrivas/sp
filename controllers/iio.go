package controllers

import (
	"net/http"
	"time"
	"path/filepath"
	"github.com/oficialrivas/sgi/utils"

	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/google/uuid"
)


// CreateIIO crea un nuevo registro de IIO
// @Summary Crea un nuevo registro de IIO
// @Description Crea un nuevo registro de IIO con los datos proporcionados
// @Tags iio
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param descripcion formData string true "Descripción del IIO"
// @Param fecha formData string true "Fecha del IIO (YYYY-MM-DD)"
// @Param lugar formData string true "Lugar del IIO"
// @Param modalidad formData string true "Modalidad del IIO"
// @Param nombre formData string true "Nombre del IIO"
// @Param parroquia formData string true "Parroquia del IIO"
// @Param redi formData string true "REDI del IIO"
// @Param zodi formData string true "ZODI del IIO"
// @Param tie formData string true "TIE del IIO"
// @Param area formData string true "Área del IIO"
// @Param imagen formData file false "Imagen del IIO"
// @Success 200 {object} models.IIO
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /iios [post]
// @Security ApiKeyAuth
func CreateIIO(c *gin.Context) {
	var input struct {
		Descripcion string `form:"descripcion" binding:"required"`
		Fecha       string `form:"fecha" binding:"required"`
		Lugar       string `form:"lugar" binding:"required"`
		Modalidad   string `form:"modalidad" binding:"required"`
		Nombre      string `form:"nombre" binding:"required"`
		Parroquia   string `form:"parroquia" binding:"required"`
		REDI        string `form:"redi" binding:"required"`
		ZODI        string `form:"zodi" binding:"required"`
		Tie         string `form:"tie" binding:"required"`
		Area        string `form:"area" binding:"required"`
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

	iio := models.IIO{
		ID:           uuid.New(),
		Descripcion:  input.Descripcion,
		Fecha:        parsedDate,
		Lugar:        input.Lugar,
		Modalidad:    input.Modalidad,
		Nombre:       input.Nombre,
		Parroquia:    input.Parroquia,
		REDI:         input.REDI,
		ZODI:         input.ZODI,
		Tie:          input.Tie,
		Area:         claims.Area, // Asignar el área del usuario desde el token
		UserID:       uuid.MustParse(userID.(string)),
		ImagenURL:    imagenURL,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := configs.DB.Create(&iio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, iio)
}

// GetIIO obtiene un IIO por su ID
// @Summary Obtiene un IIO por su ID
// @Description Obtiene los datos de un IIO por su ID
// @Tags iio
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del IIO"
// @Success 200 {object} models.IIO
// @Failure 404 {object} models.ErrorResponse
// @Router /iios/{id} [get]
// @Security ApiKeyAuth
func GetIIO(c *gin.Context) {
	id := c.Param("id")
	var iio models.IIO
	if err := configs.DB.Where("id = ?", id).First(&iio).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "IIO no encontrado"})
		return
	}

	c.JSON(http.StatusOK, iio)
}

// UpdateIIO actualiza un IIO existente
// @Summary Actualiza un IIO existente
// @Description Actualiza los datos de un IIO existente
// @Tags iio
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del IIO"
// @Param iio body models.IIO true "Datos del IIO actualizados"
// @Success 200 {object} models.IIO
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /iios/{id} [put]
// @Security ApiKeyAuth
func UpdateIIO(c *gin.Context) {
	id := c.Param("id")
	var iio models.IIO
	if err := configs.DB.Where("id = ?", id).First(&iio).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "IIO no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&iio); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	iio.ID, _ = uuid.Parse(id)
	iio.UpdatedAt = time.Now().UTC()
	if err := configs.DB.Save(&iio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, iio)
}

// DeleteIIO elimina un IIO por su ID
// @Summary Elimina un IIO por su ID
// @Description Elimina un IIO por su ID
// @Tags iio
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del IIO"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /iios/{id} [delete]
// @Security ApiKeyAuth
func DeleteIIO(c *gin.Context) {
	id := c.Param("id")
	var iio models.IIO
	if err := configs.DB.Where("id = ?", id).First(&iio).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "IIO no encontrado"})
		return
	}

	if err := configs.DB.Delete(&iio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetIIOs obtiene todas las IIOs
// @Summary Obtiene todas las IIOs
// @Description Obtiene todos los registros de IIO
// @Tags iio
// @Accept json
// @Produce json
// @Success 200 {array} models.IIO
// @Failure 500 {object} models.ErrorResponse
// @Router /iios [get]
// @Security ApiKeyAuth
func GetIIOs(c *gin.Context) {
	var iios []models.IIO
	if err := configs.DB.Find(&iios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, iios)
}


// FilterIIOs obtiene IIOs basados en filtros opcionales
// @Summary Obtiene IIOs filtrados
// @Description Obtiene registros de IIO en un periodo específico, y filtrados por tie, redi, zodi y modalidad
// @Tags iio
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Param tie query string false "TIE"
// @Param redi query string false "REDI"
// @Param zodi query string false "ZODI"
// @Param modalidad query string false "Modalidad"
// @Success 200 {array} models.IIO
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /iios/filter [get]
// @Security ApiKeyAuth
func FilterIIOs(c *gin.Context) {
	var iios []models.IIO

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tie := c.Query("tie")
	redi := c.Query("redi")
	zodi := c.Query("zodi")
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
	if modalidad != "" {
		db = db.Where("modalidad = ?", modalidad)
	}

	if err := db.Find(&iios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, iios)
}