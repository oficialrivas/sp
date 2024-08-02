package controllers

import (
	"bytes"
	"encoding/json"
	
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/utils"
)


// CreateIIO crea un nuevo registro de IIO y lo envía a un webhook
// @Summary Crea un nuevo registro de IIO y lo envía a un webhook
// @Description Crea un nuevo registro de IIO con los datos proporcionados y lo envía a un webhook
// @Tags iio
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param IIO body models.IIO true "Datos de la IIO"
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
		// Asegurar que la ruta de la imagen use barras "/"
		imagenURL = strings.ReplaceAll(imagenPath, "\\", "/")
	}

	iio := models.IIO{
		ID:           uuid.New(),
		Descripcion:  input.Descripcion,
		Fecha:        parsedDate,
		Lugar:        input.Lugar,
		Modalidad:    []models.Modalidad{}, 
		Tie:          []models.Tie{}, 
		Nombre:       input.Nombre,
		Parroquia:    input.Parroquia,
		REDI:         input.REDI,
		ZODI:         input.ZODI,
		Area:         claims.Area, 
		UserID:       uuid.MustParse(userID.(string)),
		ImagenURL:    imagenURL,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := configs.DB.Create(&iio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Enviar el mensaje al webhook
	payload := map[string]string{
		"mensaje": iio.Descripcion,
		"iio_id":  iio.ID.String(),
		"jwt":     tokenString,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal payload: %v\n", err)
	} else {
		go func() {
			resp, err := http.Post("http://10.51.16.147:8080/webhook", "application/json", bytes.NewBuffer(jsonPayload))
			if err != nil {
				fmt.Printf("Failed to send request to webhook: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Webhook responded with status: %v\n", resp.StatusCode)
			}
		}()
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



// GetIIOByParameters obtiene registros de IIO por período, REDI, temática y modalidad
// @Summary Obtiene registros de IIO por período, REDI, temática y modalidad
// @Description Recupera los registros de IIO en un período determinado, según los parámetros REDI, temática y modalidad ingresados por el usuario
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{} "result"
// @Tags iio
// @Accept json
// @Produce json
// @Param request body models.IIORequestParams true "Parámetros de consulta"
// @Success 200 {object} []models.IIO
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /gestion/iio [post]
func GetIIOByParameters(c *gin.Context) {
	var params models.IIORequestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsear las fechas de los parámetros
	startDate, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido para start_date"})
		return
	}

	endDate, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido para end_date"})
		return
	}

	// Validar que startDate no sea posterior a endDate
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date no puede ser posterior a end_date"})
		return
	}

	// Construir la consulta dinámica
	query := configs.DB.Model(&models.IIO{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)

	if params.REDI != "" {
		query = query.Where("redi = ?", params.REDI)
	}
	if params.Tie != "" {
		query = query.Where("tie = ?", params.Tie)
	}
	if params.Modalidad != "" {
		query = query.Where("modalidad = ?", params.Modalidad)
	}

	var iios []models.IIO
	if err := query.Find(&iios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver los resultados en formato JSON
	c.JSON(http.StatusOK, iios)
}


// GetIIOByModalidadAndValor obtiene registros de IIO por modalidad y valor
// @Summary Obtiene registros de IIO por modalidad y valor
// @Description Recupera los registros de IIO en un período determinado, según los parámetros modalidad y valor ingresados por el usuario
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Tags iio
// @Accept json
// @Produce json
// @Param request body models.IIORequestParams true "Parámetros de consulta"
// @Success 200 {object} []models.IIO
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /gestion/iio/modalidad [post]
func GetIIOByModalidadAndValor(c *gin.Context) {
	var params models.IIORequestParams2
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsear las fechas de los parámetros
	startDate, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido para start_date"})
		return
	}

	endDate, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido para end_date"})
		return
	}

	// Validar que startDate no sea posterior a endDate
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date no puede ser posterior a end_date"})
		return
	}

	// Construir la consulta dinámica
	query := configs.DB.Model(&models.IIO{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)

	if params.Modalidad != "" {
		query = query.Where("modalidad = ?", params.Modalidad)
	}
	if params.Valor != nil {
		query = query.Where("valor = ?", *params.Valor)
	}

	var iios []models.IIO
	if err := query.Find(&iios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver los resultados en formato JSON
	c.JSON(http.StatusOK, iios)
}

// GetIIOCountByModalidadAndValor obtiene el número de registros de IIO por modalidad y valor
// @Summary Obtiene el número de registros de IIO por modalidad y valor
// @Description Recupera el número de registros de IIO en un período determinado, según los parámetros modalidad y valor ingresados por el usuario
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Tags iio
// @Accept json
// @Produce json
// @Param request body models.IIORequestParams3 true "Parámetros de consulta"
// @Success 200 {object} map[string]interface{} "result"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /gestion/iio/modalidad/count [post]
func GetIIOCountByModalidadAndValor(c *gin.Context) {
	var params models.IIORequestParams3
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsear las fechas de los parámetros
	startDate, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido para start_date"})
		return
	}

	endDate, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido para end_date"})
		return
	}

	// Validar que startDate no sea posterior a endDate
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date no puede ser posterior a end_date"})
		return
	}

	// Crear una subconsulta para obtener los IDs de IIOs con la modalidad especificada
	var iioIDs []uuid.UUID
	if err := configs.DB.Table("iio_modalidad").
		Select("iio_id").
		Joins("JOIN modalidad ON modalidad.id = iio_modalidad.modalidad_id").
		Where("modalidad.nombre = ?", params.Modalidad).
		Pluck("iio_id", &iioIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(iioIDs) == 0 {
		// Si no hay IDs de IIOs con la modalidad especificada, devolver conteos de 0
		c.JSON(http.StatusOK, map[string]interface{}{
			"count_true":  0,
			"count_false": 0,
		})
		return
	}

	fmt.Printf("IIO IDs: %v\n", iioIDs)

	// Contar registros con Valor en true
	var countTrue int64
	if err := configs.DB.Model(&models.IIO{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("id IN (?)", iioIDs).
		Where("valor = ?", true).
		Count(&countTrue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Contar registros con Valor en false
	var countFalse int64
	if err := configs.DB.Model(&models.IIO{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("id IN (?)", iioIDs).
		Where("valor = ?", false).
		Count(&countFalse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver los resultados en formato JSON
	result := map[string]interface{}{
		"count_true":  countTrue,
		"count_false": countFalse,
	}

	c.JSON(http.StatusOK, result)
}