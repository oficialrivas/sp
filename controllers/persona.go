package controllers

import (
	"net/http"
	"time"
	"log"
	"strings"
	"github.com/oficialrivas/sgi/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	configs "github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
)

// CreatePersona crea un nuevo registro de Persona
// @Summary Crea un nuevo registro de Persona
// @Description Crea un nuevo registro de Persona con los datos proporcionados
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param persona body models.Persona true "Datos del Persona"
// @Success 200 {object} models.Persona
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /personas [post]
// @Security ApiKeyAuth
func CreatePersona(c *gin.Context) {
	var persona models.Persona
	if err := c.ShouldBindJSON(&persona); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	persona.UserID = uuid.MustParse(userID.(string))

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
	persona.Area = claims.Area

	if err := configs.DB.Create(&persona).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, persona)
}

// GetPersona obtiene un Persona por su ID
// @Summary Obtiene un Persona por su ID
// @Description Obtiene los datos de un Persona por su ID
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del Persona"
// @Success 200 {object} models.Persona
// @Failure 404 {object} models.ErrorResponse
// @Router /personas/{id} [get]
// @Security ApiKeyAuth
func GetPersona(c *gin.Context) {
	id := c.Param("id")
	var persona models.Persona
	if err := configs.DB.Where("id = ?", id).Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").First(&persona).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Persona no encontrado"})
		return
	}

	c.JSON(http.StatusOK, persona)
}

// UpdatePersona actualiza un Persona existente
// @Summary Actualiza un Persona existente
// @Description Actualiza los datos de un Persona existente
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del Persona"
// @Param persona body models.Persona true "Datos del Persona actualizados"
// @Success 200 {object} models.Persona
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /personas/{id} [put]
// @Security ApiKeyAuth
func UpdatePersona(c *gin.Context) {
	id := c.Param("id")
	var persona models.Persona
	if err := configs.DB.Where("id = ?", id).First(&persona).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Persona no encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&persona); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	persona.ID, _ = uuid.Parse(id)
	persona.UpdatedAt = time.Now().UTC()
	if err := configs.DB.Save(&persona).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, persona)
}

// DeletePersona elimina un Persona por su ID
// @Summary Elimina un Persona por su ID
// @Description Elimina un Persona por su ID
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del Persona"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /personas/{id} [delete]
// @Security ApiKeyAuth
func DeletePersona(c *gin.Context) {
	id := c.Param("id")
	var persona models.Persona
	if err := configs.DB.Where("id = ?", id).First(&persona).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Persona no encontrado"})
		return
	}

	if err := configs.DB.Delete(&persona).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPersonaByCedula obtiene un Persona por su número de cédula
// @Summary Obtiene un Persona por su número de cédula
// @Description Obtiene los datos de un Persona por su número de cédula
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param cedula path string true "Cédula del Persona"
// @Success 200 {object} models.Persona
// @Failure 404 {object} models.ErrorResponse
// @Router /personas/cedula/{cedula} [get]
// @Security ApiKeyAuth
func GetPersonaByCedula(c *gin.Context) {
	cedula := c.Param("cedula")
	var persona models.Persona
	if err := configs.DB.Where("cedula = ?", cedula).Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").First(&persona).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Persona no encontrado"})
		return
	}

	c.JSON(http.StatusOK, persona)
}

// GetPersonaByPasaporte obtiene un Persona por su número de pasaporte
// @Summary Obtiene un Persona por su número de pasaporte
// @Description Obtiene los datos de un Persona por su número de pasaporte
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param pasaporte path string true "Pasaporte del Persona"
// @Success 200 {object} models.Persona
// @Failure 404 {object} models.ErrorResponse
// @Router /personas/pasaporte/{pasaporte} [get]
// @Security ApiKeyAuth
func GetPersonaByPasaporte(c *gin.Context) {
	pasaporte := c.Param("pasaporte")
	var persona models.Persona
	if err := configs.DB.Where("pasaporte = ?", pasaporte).Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").First(&persona).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Persona no encontrado"})
		return
	}

	c.JSON(http.StatusOK, persona)
}

// GetPersonaByNombre obtiene un Persona por su nombre
// @Summary Obtiene un Persona por su nombre
// @Description Obtiene los datos de un Persona por su nombre
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param nombre path string true "Nombre del Persona"
// @Success 200 {object} models.Persona
// @Failure 404 {object} models.ErrorResponse
// @Router /personas/nombre/{nombre} [get]
// @Security ApiKeyAuth
func GetPersonaByNombre(c *gin.Context) {
	nombre := c.Param("nombre")
	var persona models.Persona
	if err := configs.DB.Where("nombre = ?", nombre).Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").First(&persona).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Persona no encontrado"})
		return
	}

	c.JSON(http.StatusOK, persona)
}

// GetPersonasByNacionalidad obtiene todos los Personas por su nacionalidad
// @Summary Obtiene todos los Personas por su nacionalidad
// @Description Obtiene los datos de todos los Personas que pertenecen a una nacionalidad específica
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param nacionalidad path string true "Nacionalidad de los Personas"
// @Success 200 {array} models.Persona
// @Failure 404 {object} models.ErrorResponse
// @Router /personas/nacionalidad/{nacionalidad} [get]
// @Security ApiKeyAuth
func GetPersonasByNacionalidad(c *gin.Context) {
	nacionalidad := c.Param("nacionalidad")
	var personas []models.Persona
	if err := configs.DB.Where("nacionalidad = ?", nacionalidad).Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").Find(&personas).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Personas no encontrados"})
		return
	}

	c.JSON(http.StatusOK, personas)
}

// GetPersonas obtiene todos los Personas
// @Summary Obtiene todos los Personas
// @Description Obtiene los datos de todos los Personas
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Persona
// @Failure 404 {object} models.ErrorResponse
// @Router /personas [get]
// @Security ApiKeyAuth
func GetPersonas(c *gin.Context) {
	var personas []models.Persona
	if err := configs.DB.Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").Find(&personas).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Personas no encontrados"})
		return
	}

	c.JSON(http.StatusOK, personas)
}

// GetPersonasByCedula obtiene personas por una lista de cédulas
// @Summary Obtiene personas por una lista de cédulas
// @Description Obtiene los datos de personas por una lista de cédulas
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param cedulas query string true "Lista de cédulas separadas por comas"
// @Success 200 {array} models.Persona
// @Failure 404 {object} models.ErrorResponse
// @Router /personas/cedulas [get]
// @Security ApiKeyAuth
func GetPersonasByCedula(c *gin.Context) {
	cedulas := c.Query("cedulas")
	if cedulas == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "No se proporcionó ninguna cédula"})
		return
	}

	cedulaList := strings.Split(cedulas, ",")
	var personas []models.Persona
	if err := configs.DB.Where("cedula IN (?)", cedulaList).Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").Find(&personas).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Personas no encontradas"})
		return
	}

	c.JSON(http.StatusOK, personas)
}

// SearchPersonas busca personas usando texto completo
// @Summary Busca personas usando texto completo
// @Description Busca personas en la tabla persona usando un índice de texto completo
// @Tags persona
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param query query string true "Texto de búsqueda"
// @Success 200 {array} models.Persona
// @Failure 400 {object} models.ErrorResponse
// @Router /personas/search [get]
func SearchPersonas(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "El parámetro de búsqueda es obligatorio"})
		return
	}

	// Convertir query a una consulta tsquery adecuada
	tsQuery := strings.Join(strings.Fields(query), " & ")

	// Log para depuración
	log.Printf("Running fulltext search with query: %s", tsQuery)

	var personas []models.Persona
	err := configs.DB.Raw(`
		SELECT * FROM persona
		WHERE to_tsvector('spanish', coalesce(nombre, '') || ' ' || coalesce(apellido, '') || ' ' || coalesce(telefono, '') || ' ' || coalesce(profesion, '') || ' ' || coalesce(cedula, '') || ' ' || coalesce(correo, '')) @@ to_tsquery('spanish', ?)
	`, tsQuery).Scan(&personas).Error
	if err != nil {
		log.Printf("Error en la búsqueda: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Error en la búsqueda"})
		return
	}

	// Log resultado de la búsqueda
	log.Printf("Found %d personas", len(personas))

	c.JSON(http.StatusOK, personas)
}
