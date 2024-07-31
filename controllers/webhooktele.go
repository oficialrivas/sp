package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
)

type GenerateTokenRequest2 struct {
	Telefono string `json:"telefono"`
}

type TokenResponse2 struct {
	AccessToken  string `json:"accessToken"` // Asegúrate de que el nombre del campo coincida
	RefreshToken string `json:"refreshToken"`
	ID           string `json:"id"`
}

type UpdateTelegramRequest struct {
	Usuario string `json:"u_telegram" binding:"required"`
}

func Websmstelegram(c *gin.Context) {
	var data struct {
		Canal    string `json:"canal"`
		APIKey   string `json:"api_key"`
		Numero   string `json:"numero"`
		Usuario  struct {
			ID        string `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"usuario"`
		PhotoURL string `json:"photo_url"`
		Text     string `json:"text"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON", "details": err.Error()})
		return
	}
	fmt.Printf("Received data: %+v\n", data)

	// Procesar el mensaje y verificar el patrón PC(04127100820)
	re := regexp.MustCompile(`PC\((\d+)\)`)
	matches := re.FindStringSubmatch(data.Text)
	if len(matches) == 2 {
		telefono := matches[1]
		usuarioID := data.Usuario.ID

		// Verificar que el ID de usuario no esté vacío
		if usuarioID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "u_telegram (Telegram ID) is empty"})
			return
		}

		// Generar el token JWT utilizando el número de teléfono
		payload := GenerateTokenRequest2{
			Telefono: telefono,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload", "details": err.Error()})
			return
		}
		fmt.Printf("Sending payload for token generation: %s\n", string(jsonPayload))

		resp, err := http.Post("http://localhost:8080/generate-token", "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Leer y devolver la respuesta
		var tokenResponse TokenResponse2
		if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode token response", "details": err.Error()})
			return
		}
		fmt.Printf("Token generation response: %+v\n", tokenResponse)

		// Verificar que el token JWT haya sido generado correctamente
		if tokenResponse.AccessToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed, access token is empty"})
			return
		}

		// Datos para el endpoint de actualización
		url := fmt.Sprintf("http://localhost:8080/users/telefono/%s", telefono)
		headers := map[string]string{
			"Authorization": tokenResponse.AccessToken, // Usar el token de acceso
			"Content-Type":  "application/json",
		}
		updatePayload := UpdateTelegramRequest{
			Usuario: usuarioID,
		}
		updateJsonPayload, err := json.Marshal(updatePayload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal update payload", "details": err.Error()})
			return
		}
		fmt.Printf("Sending payload for updating u_telegram: %s\n", string(updateJsonPayload))

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(updateJsonPayload))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request", "details": err.Error()})
			return
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request", "details": err.Error()})
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body", "details": err.Error()})
			return
		}
		bodyString := string(bodyBytes)

		fmt.Printf("Response from update u_telegram: %s\n", bodyString)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to update u_telegram: %s\n", bodyString)
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to update u_telegram", "details": bodyString})
			return
		}

		fmt.Println("u_telegram actualizado correctamente.")

		// Fetch the user details using the ID from token response
		fmt.Printf("Fetching user details for ID: %s\n", tokenResponse.ID)
		var user models.User
		if err := configs.DB.First(&user, "id = ?", tokenResponse.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "details": err.Error()})
			return
		}

		// Download the image
		imagePath := ""
		if data.PhotoURL != "" {
			resp, err := http.Get(data.PhotoURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download image", "details": err.Error()})
				return
			}
			defer resp.Body.Close()

			// Create the static directory if it doesn't exist
			if _, err := os.Stat("static"); os.IsNotExist(err) {
				os.Mkdir("static", os.ModePerm)
			}

			// Create the file
			imageFileName := uuid.New().String() + filepath.Ext(data.PhotoURL)
			imagePath = filepath.Join("static", imageFileName)
			file, err := os.Create(imagePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image file", "details": err.Error()})
				return
			}
			defer file.Close()

			// Copy the downloaded image to the file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image", "details": err.Error()})
				return
			}

			// Ensure the image path uses forward slashes
			imagePath = strings.ReplaceAll(imagePath, "\\", "/")
		}

		// Create the message record
		parsedDate := time.Now() // Using current time as message date for example

		mensaje := models.Mensaje{
			Descripcion: data.Text,
			Fecha:       parsedDate,
			Lugar:       "", // Assuming this comes from elsewhere, set as needed
			Modalidad:   "", // Assuming this comes from elsewhere, set as needed
			Nombre:      data.Usuario.FirstName + " " + data.Usuario.LastName, // Combining first and last name
			Parroquia:   "", // Assuming this comes from elsewhere, set as needed
			Canal:       "Telegram", // Assuming this comes from elsewhere, set as needed
			REDI:        user.REDI,
			ZODI:        user.Zodi,
			ADI:         user.ADI,
			Tie:         "", // Assuming this comes from elsewhere, set as needed
			Area:        user.Area, // Assuming this comes from JWT claims
			UserID:      uuid.MustParse(tokenResponse.ID),
			ImagenURL:   imagePath, // Use the path of the saved image
			Procesado:   false,     // Assuming this comes from elsewhere, set as needed
		}

		if err := configs.DB.Create(&mensaje).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, mensaje)
	} 
}
