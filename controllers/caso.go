package controllers

import (
    "net/http"
    "time"
	"gorm.io/gorm"
    "log"
	

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/oficialrivas/sgi/config"
    "github.com/oficialrivas/sgi/models"
)

// CreateCaso crea un nuevo caso
// @Summary Crea un nuevo caso
// @Accept json
// @Produce json
// @Tags Caso
// @Param Authorization header string true "Bearer token"
// @Param caso body models.Caso true "Datos del caso"
// @Success 201 {object} models.Caso
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /casos [post]
func CreateCaso(c *gin.Context) {
    var caso models.Caso
    if err := c.ShouldBindJSON(&caso); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
        return
    }

  

    userUUID, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
        return
    }

    // Verificar si el usuario existe en la base de datos
    var user models.User
    if err := configs.DB.First(&user, "id = ?", userUUID).Error; err != nil {
        log.Printf("User not found: %s", userUUID) // Log user not found
        c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
    }

    // Añadir el userID al caso
    caso.UserID = user.ID

    caso.ID = uuid.New()
    caso.CreatedAt = time.Now().UTC()
    caso.UpdatedAt = caso.CreatedAt

    if err := configs.DB.Create(&caso).Error; err != nil {
        log.Printf("Error creating case: %v", err) // Log error creating case
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el caso"})
        return
    }

    c.JSON(http.StatusCreated, caso)
}

// GetCasoByID obtiene un caso por su ID
// @Summary Obtiene un caso por su ID
// @Produce json
// @Tags Caso
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del caso"
// @Success 200 {object} models.Caso
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /casos/{id} [get]
func GetCasoByID(c *gin.Context) {
	id := c.Param("id")

	var caso models.Caso
	if err := configs.DB.Preload("Relacion").Preload("Vehiculos").Preload("Empresas").Preload("Direcciones").Preload("IIOs").Preload("Documentos").First(&caso, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Caso no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el caso"})
		return
	}

	c.JSON(http.StatusOK, caso)
}

// UpdateCaso actualiza un caso existente por su ID
// @Summary Actualiza un caso existente por su ID
// @Accept json
// @Produce json
// @Tags Caso
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del caso a actualizar"
// @Param caso body models.Caso true "Datos del caso a actualizar"
// @Success 200 {object} models.Caso
// @Failure 400 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /casos/{id} [put]
func UpdateCaso(c *gin.Context) {
	id := c.Param("id")

	var caso models.Caso
	if err := configs.DB.First(&caso, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Caso no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el caso"})
		return
	}

	if err := c.ShouldBindJSON(&caso); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	caso.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&caso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el caso"})
		return
	}

	c.JSON(http.StatusOK, caso)
}

// DeleteCaso elimina un caso por su ID
// @Summary Elimina un caso por su ID
// @Produce json
// @Tags Caso
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del caso a eliminar"
// @Success 204 "Caso eliminado exitosamente"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /casos/{id} [delete]
func DeleteCaso(c *gin.Context) {
	id := c.Param("id")

	var caso models.Caso
	if err := configs.DB.First(&caso, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Caso no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el caso"})
		return
	}

	if err := configs.DB.Delete(&caso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el caso"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ValorarCaso permite a los usuarios valorar un caso según su rol
// @Summary Valorar un caso
// @Accept json
// @Produce json
// @Tags Caso
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del caso"
// @Param valor body int true "Valoración"
// @Success 200 {object} models.Caso
// @Failure 400 {object} map[string]string "error"
// @Failure 403 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /casos/valorar/{id} [put]
func ValorarCaso(c *gin.Context) {
    var input struct {
        Valor int `json:"valor" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    casoID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
        return
    }

    userID, _ := c.Get("userID")
    role, _ := c.Get("role")

    var caso models.Caso
    if err := configs.DB.First(&caso, "id = ?", casoID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    switch role {
    case "admin":
        caso.Vdirector = input.Valor
    case "analista":
        caso.Vanalista = input.Valor
    case "superusuario":
        caso.Vcoordinador = input.Valor
    default:
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized role"})
        return
    }

    // Añadir el usuario al caso
    userUUID, _ := uuid.Parse(userID.(string))
    user := models.User{ID: userUUID}
    caso.Users = append(caso.Users, user)

    if err := configs.DB.Save(&caso).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el caso"})
        return
    }

    c.JSON(http.StatusOK, caso)
}