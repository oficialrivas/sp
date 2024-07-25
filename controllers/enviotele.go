package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/utils"
)

// sendMessageInput es la estructura para el cuerpo de la solicitud de SendMessageHandler
type sendMessageInput struct {
	Mensaje string `json:"mensaje" binding:"required"`
}

// SendMessageHandler recibe el JSON con el id del usuario y el mensaje, y envía el mensaje a través del webhook
// @Summary Envía un mensaje a un usuario específico basado en su ID
// @Description Envía un mensaje a un usuario específico basado en su ID
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del usuario"
// @Param mensaje body sendMessageInput true "Mensaje"
// @Success 200 {object} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /send_telegram/{id} [post]
// @Security BearerAuth
func SendMessageHandler(c *gin.Context) {
	var input sendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del usuario desde los parámetros de la URL
	userIDParam := c.Param("id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
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

	authUserID := claims.UserID
	if authUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	// Obtener los detalles del usuario usando el ID del parámetro
	var user models.User
	if err := configs.DB.First(&user, "id = ?", userIDParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Crear el registro del mensaje
	mensaje := models.Mensaje{
		Descripcion: input.Mensaje,
		Fecha:       time.Now(), // Establecer la fecha automáticamente
		REDI:        user.REDI,
		ZODI:        user.Zodi,
		ADI:         user.ADI,
		Tie:         "", // Assuming this comes from elsewhere, set as needed
		UserID:      user.ID,
		Canal:        "Telegram", // Establecer el tipo como "saliente"
		Tipo:        "saliente", // Establecer el tipo como "saliente"
		Procesado:   false, // Assuming this should be set to false when sending the message
	}

	if err := configs.DB.Create(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preparar la carga útil para enviar al webhook
	payload := map[string]interface{}{
		"id":      user.Usuario,
		"mensaje": input.Mensaje,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload", "details": err.Error()})
		return
	}

	// Enviar el mensaje al webhook
	resp, err := http.Post("http://localhost:5000/send_message", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to webhook", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Webhook responded with status: %v", resp.StatusCode)})
		return
	}

	c.JSON(http.StatusOK, mensaje)
}



// SendMessagesByRediHandler envía mensajes a todos los usuarios en una REDI específica
// @Summary Envía mensajes a todos los usuarios en una REDI específica
// @Description Envía mensajes a todos los usuarios en una REDI específica
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param redi path string true "REDI del usuario"
// @Param mensaje body sendMessageInput true "Mensaje"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /send_messages_by_redi_tele/{redi} [post]
// @Security BearerAuth
func SendMessagesByRediHandler(c *gin.Context) {
	var input sendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener la REDI desde los parámetros de la URL
	rediParam := c.Param("redi")
	if rediParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "REDI is required"})
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

	authUserID := claims.UserID
	if authUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	// Obtener los detalles de los usuarios por REDI
	var users []models.User
	if err := configs.DB.Where("redi = ?", rediParam).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users", "details": err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found in the specified REDI"})
		return
	}

	// Crear y enviar mensajes a cada usuario
	var messages []models.Mensaje
	for _, user := range users {
		// Crear el registro del mensaje
		mensaje := models.Mensaje{
			Descripcion: input.Mensaje,
			Fecha:       time.Now(), // Establecer la fecha automáticamente
			REDI:        user.REDI,
			ZODI:        user.Zodi,
			ADI:         user.ADI,
			UserID:      user.ID,
			Canal:       "Telegram",
			Tipo:        "saliente",
			Procesado:   false,
		}

		if err := configs.DB.Create(&mensaje).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Preparar la carga útil para enviar al webhook
		payload := map[string]interface{}{
			"id":      user.Usuario,
			"mensaje": input.Mensaje,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload", "details": err.Error()})
			return
		}

		// Enviar el mensaje al webhook
		resp, err := http.Post("http://localhost:5000/send_message", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to webhook", "details": err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Webhook responded with status: %v", resp.StatusCode)})
			return
		}

		messages = append(messages, mensaje)
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}