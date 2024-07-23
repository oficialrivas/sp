package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/oficialrivas/sgi/config"
    "github.com/oficialrivas/sgi/models"
    "github.com/google/uuid"
)

// UserPeriodRequestParams es la estructura para los parámetros del request
type UserPeriodRequestParams struct {
    UserID    string `json:"user_id" binding:"required"`
    StartDate string `json:"start_date" binding:"required"`
    EndDate   string `json:"end_date" binding:"required"`
}

// RecordsCountByArea contiene los conteos de registros por área
type UserRecordsCountByArea struct {
    Area             string `json:"area"`
    CasosCount       int64  `json:"casos_count"`
    PersonasCount    int64  `json:"personas_count"`
    VehiculosCount   int64  `json:"vehiculos_count"`
    EmpresasCount    int64  `json:"empresas_count"`
    DireccionesCount int64  `json:"direcciones_count"`
    IIOsCount        int64  `json:"iios_count"`
    DocumentosCount  int64  `json:"documentos_count"`
}

// ModalidadCount contiene los conteos de registros por modalidad
type ModalidadCount1 struct {
    Modalidad string `json:"modalidad"`
    Count     int64  `json:"count"`
}

// AreaModalidadCount contiene los conteos de registros por modalidad y área
type UserAreaModalidadCount struct {
    Area  string          `json:"area"`
    Casos []ModalidadCount `json:"casos"`
    IIOs  []ModalidadCount `json:"iios"`
}

// GetRecordsByUserAndPeriod devuelve el número de registros de un área en un período específico para un usuario específico
// @Summary Devuelve el número de registros de un área en un período específico para un usuario específico
// @Description Obtiene el número de registros para un área específica en un período determinado para un usuario específico
// @Tags Registros
// @Accept json
// @Produce json
// @Param request body UserPeriodRequestParams true "Parámetros de consulta"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /gestion/user [post]
func GetRecordsByUserAndPeriod(c *gin.Context) {
    var params UserPeriodRequestParams
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

    // Contadores para cada modelo
    var casosCount, personasCount, vehiculosCount, empresasCount, direccionesCount, iiosCount, documentosCount int64

    userUUID, err := uuid.Parse(params.UserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de UUID inválido para user_id"})
        return
    }

    // Consultar en cada modelo
    configs.DB.Model(&models.Caso{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userUUID, startDate, endDate).Count(&casosCount)
    configs.DB.Model(&models.Persona{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userUUID, startDate, endDate).Count(&personasCount)
    configs.DB.Model(&models.Vehiculo{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userUUID, startDate, endDate).Count(&vehiculosCount)
    configs.DB.Model(&models.Empresa{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userUUID, startDate, endDate).Count(&empresasCount)
    configs.DB.Model(&models.Direccion{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userUUID, startDate, endDate).Count(&direccionesCount)
    configs.DB.Model(&models.IIO{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userUUID, startDate, endDate).Count(&iiosCount)
    configs.DB.Model(&models.Documento{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userUUID, startDate, endDate).Count(&documentosCount)

    // Devolver los resultados en formato JSON
    c.JSON(http.StatusOK, gin.H{
        "casos_count":       casosCount,
        "personas_count":    personasCount,
        "vehiculos_count":   vehiculosCount,
        "empresas_count":    empresasCount,
        "direcciones_count": direccionesCount,
        "iios_count":        iiosCount,
        "documentos_count":  documentosCount,
    })
}

// GetRecordsCountByUserAndModalidad devuelve el número de registros por área y modalidad en un período específico para un usuario específico
// @Summary Devuelve el número de registros por área y modalidad en un período específico para un usuario específico
// @Description Obtiene el número de registros por área y modalidad en un período determinado para un usuario específico
// @Tags Registros
// @Accept json
// @Produce json
// @Param request body UserPeriodRequestParams true "Parámetros de consulta"
// @Success 200 {object} []UserAreaModalidadCount
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /gestion/user-area-modalidad [post]
func GetRecordsCountByUserAndModalidad(c *gin.Context) {
    var params UserPeriodRequestParams
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

    userUUID, err := uuid.Parse(params.UserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de UUID inválido para user_id"})
        return
    }

    // Lista para almacenar los resultados por área
    var results []UserAreaModalidadCount

    // Obtener todas las áreas únicas de la tabla Caso y IIO
    var areas []string
    configs.DB.Model(&models.Caso{}).Distinct("area").Pluck("area", &areas)
    configs.DB.Model(&models.IIO{}).Distinct("area").Pluck("area", &areas)

    // Iterar sobre cada área y obtener los conteos por modalidad
    for _, area := range areas {
        var casos []ModalidadCount
        var iios []ModalidadCount

        configs.DB.Model(&models.Caso{}).
            Select("modalidad, COUNT(*) as count").
            Where("user_id = ? AND area = ? AND created_at BETWEEN ? AND ?", userUUID, area, startDate, endDate).
            Group("modalidad").
            Find(&casos)

        configs.DB.Model(&models.IIO{}).
            Select("modalidad, COUNT(*) as count").
            Where("user_id = ? AND area = ? AND created_at BETWEEN ? AND ?", userUUID, area, startDate, endDate).
            Group("modalidad").
            Find(&iios)

        result := UserAreaModalidadCount{
            Area:  area,
            Casos: casos,
            IIOs:  iios,
        }

        results = append(results, result)
    }

    // Devolver los resultados en formato JSON
    c.JSON(http.StatusOK, results)
}
