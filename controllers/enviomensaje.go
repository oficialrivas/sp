package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/utils"
)

// CreateAndSendMensaje crea un mensaje y lo envía a un webhook
// @Summary Crea un nuevo registro de Mensaje y lo envía a un webhook
// @Description Crea un nuevo registro de Mensaje con los datos proporcionados y lo envía a un webhook
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param mensaje body createAndSendMensajeInput true "Mensaje"
// @Success 200 {object} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /create-and-send-mensaje [post]
// @Security BearerAuth
func CreateAndSendMensaje(c *gin.Context) {
	var input createAndSendMensajeInput
	if err := c.ShouldBindJSON(&input); err != nil {
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

	userID := claims.UserID
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	// Crear el registro del mensaje
	mensaje := models.Mensaje{
		Descripcion: input.Descripcion,
		Fecha:       time.Now(), // Establecer la fecha automáticamente
		REDI:        input.REDI,
		ZODI:        input.ZODI,
		ADI:         input.ADI,
		Tie:         input.Tie,
		UserID:      uuid.MustParse(userID),
		Tipo:        "saliente", // Establecer el tipo como "saliente"
		Procesado:   input.Procesado,
	}

	if err := configs.DB.Create(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener los usuarios que coincidan con los parámetros redi, zodi, adi
	var users []models.User
	if err := configs.DB.Where("redi = ? OR zodi = ? OR adi = ?", input.REDI, input.ZODI, input.ADI).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enviar el mensaje a cada usuario encontrado
	for _, user := range users {
		go sendToWebhook(user.Telefono, input.Descripcion, input.REDI)
	}

	c.JSON(http.StatusOK, mensaje)
}

// SendMensajeToUser envía un mensaje a un usuario específico basado en su ID
// @Summary Envía un mensaje a un usuario específico basado en su ID
// @Description Envía un mensaje a un usuario específico basado en su ID
// @Tags mensaje
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_id path string true "ID del usuario"
// @Param mensaje body sendMensajeToUserInput true "Mensaje"
// @Success 200 {object} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /send-mensaje-to-user/{user_id} [post]
// @Security BearerAuth
func SendMensajeToUser(c *gin.Context) {
	var input sendMensajeToUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID del usuario desde los parámetros de la URL
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Obtener el token desde el encabezado
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not found"})
		return
	}

	// Validar el token
	if _, err := utils.ValidateJWT(tokenString, false); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
		return
	}

	// Obtener los detalles del usuario
	var user models.User
	if err := configs.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Crear el registro del mensaje
	mensaje := models.Mensaje{
		Descripcion: input.Descripcion,
		Fecha:       time.Now(), // Establecer la fecha automáticamente
		REDI:        user.REDI,
		ZODI:        user.Zodi,
		ADI:         user.ADI,
		Tie:         input.Tie,
		UserID:      userID,
		Tipo:        "saliente", // Establecer el tipo como "saliente"
		Procesado:   false, // Assuming this should be set to false when sending the message
	}

	if err := configs.DB.Create(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enviar el mensaje al usuario específico
	sendToWebhook(user.Telefono, input.Descripcion, user.REDI)

	c.JSON(http.StatusOK, mensaje)
}

// createAndSendMensajeInput es la estructura para el cuerpo de la solicitud de CreateAndSendMensaje
type createAndSendMensajeInput struct {
	Descripcion string `json:"descripcion" binding:"required"`
	REDI        string `json:"redi" binding:"required"`
	ZODI        string `json:"zodi" binding:"required"`
	ADI         string `json:"adi" binding:"required"`
	Tie         string `json:"tie" binding:"required"`
	Procesado   bool   `json:"procesado" binding:"required"`
}

// sendMensajeToUserInput es la estructura para el cuerpo de la solicitud de SendMensajeToUser
type sendMensajeToUserInput struct {
	Descripcion string `json:"descripcion" binding:"required"`
	Tie         string `json:"tie" binding:"required"`
}

// sendToWebhook envía una solicitud POST al webhook con el mensaje y los detalles del usuario
func sendToWebhook(number string, text string, redi string) {
	smsc := map[string]string{
		"Guayana": "gua-llan", "los llanos": "gua-llan",
		"central": "cen-cap", "capital": "cen-cap",
		"andes": "and-occi", "occidental": "and-occi",
		"maritima": "main-ori", "oriental": "main-ori",
	}

	payload := map[string]interface{}{
		"Text":   text,
		"SMSC":   map[string]interface{}{"Location": smsc[redi]},
		"Number": number,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal payload: %v\n", err)
		return
	}

	resp, err := http.Post("http://webhook.url", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Failed to send request to webhook: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Webhook responded with status: %v\n", resp.StatusCode)
	}
}
