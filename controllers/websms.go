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
	
)

type GenerateTokenRequest struct {
	Telefono string `json:"telefono"`
	Usuario  string `json:"usuario"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ID           string `json:"id"`
}

// WebsmsHandler maneja las solicitudes entrantes al webhook y las muestra en la terminal
func WebsmsHandler(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Received data: %+v\n", data)

	senderNumber, ok := data["SenderNumber"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SenderNumber not found"})
		return
	}

	textDecoded, ok := data["TextDecoded"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TextDecoded not found"})
		return
	}

	// Prepare the request payload
	payload := map[string]string{"telefono": senderNumber}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload"})
		return
	}

	

	// Make a POST request to /generate-token
	resp, err := http.Post("http://localhost:8080/generate-token", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	defer resp.Body.Close()

	// Read and return the response
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode token response"})
		return
	}

	// Fetch the user details using the ID from token response
	var user models.User
	if err := configs.DB.First(&user, "id = ?", tokenResponse.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Create the message record
	parsedDate := time.Now() // Using current time as message date for example

	mensaje := models.Mensaje{
		Descripcion: textDecoded,
		Fecha:       parsedDate,
		Lugar:       "", // Assuming this comes from elsewhere, set as needed
		Modalidad:   "", // Assuming this comes from elsewhere, set as needed
		Nombre:      user.Nombre + " " + user.Apellido, // Combining first and last name
		Parroquia:   "", // Assuming this comes from elsewhere, set as needed
		Canal:   "SMS", // Assuming this comes from elsewhere, set as needed
		REDI:        user.REDI,
		ZODI:        user.Zodi,
		ADI:         user.ADI,
		Tie:         "", // Assuming this comes from elsewhere, set as needed
		Area:        user.Area, // Assuming this comes from JWT claims
		UserID:      uuid.MustParse(tokenResponse.ID),
		ImagenURL:   "", // Assuming this comes from elsewhere, set as needed
		Procesado:   false, // Assuming this comes from elsewhere, set as needed
	}

	if err := configs.DB.Create(&mensaje).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mensaje)
}
