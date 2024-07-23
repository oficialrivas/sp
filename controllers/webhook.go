package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/models"
)

type LoginRequest struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ID           string `json:"id"`
}

func getJWTToken(loginURL string, loginData LoginRequest) (string, error) {
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return "", err
	}

	return loginResponse.AccessToken, nil
}

func createOrUpdatePersona(persona models.Persona, token string) error {
	jsonData, err := json.Marshal(persona)
	if err != nil {
		return err
	}

	client := &http.Client{}

	// Consultar la persona por cédula
	log.Printf("Sending request to check if persona exists with Cedula: %s", persona.Cedula)
	req, err := http.NewRequest("GET", "http://localhost:8080/personas/cedula/"+persona.Cedula, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Crear nueva persona
		log.Printf("Persona not found, creating new persona with Cedula: %s", persona.Cedula)
		req, err = http.NewRequest("POST", "http://localhost:8080/personas", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("create persona request failed with status %d", resp.StatusCode)
		}
		log.Printf("Persona created successfully with Cedula: %s", persona.Cedula)
	} else if resp.StatusCode == http.StatusOK {
		// Actualizar persona existente
		log.Printf("Persona found, updating existing persona with Cedula: %s", persona.Cedula)
		req, err = http.NewRequest("PUT", "http://localhost:8080/personas/"+persona.ID.String(), bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("update persona request failed with status %d", resp.StatusCode)
		}
		log.Printf("Persona updated successfully with Cedula: %s", persona.Cedula)
	} else {
		log.Printf("Persona already exists with Cedula: %s", persona.Cedula)
	}

	return nil
}

func createOrUpdateVehiculo(vehiculo models.Vehiculo, token string) error {
	jsonData, err := json.Marshal(vehiculo)
	if err != nil {
		return err
	}

	client := &http.Client{}

	// Consultar el vehículo por matrícula
	log.Printf("Sending request to check if vehiculo exists with Matricula: %s", vehiculo.Matricula)
	req, err := http.NewRequest("GET", "http://localhost:8080/vehiculos/search?matricula="+vehiculo.Matricula, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Crear nuevo vehículo
		log.Printf("Vehiculo not found, creating new vehiculo with Matricula: %s", vehiculo.Matricula)
		req, err = http.NewRequest("POST", "http://localhost:8080/vehiculos", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("create vehiculo request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}
		log.Printf("Vehiculo created successfully with Matricula: %s", vehiculo.Matricula)
	} else if resp.StatusCode == http.StatusOK {
		// Actualizar vehículo existente
		log.Printf("Vehiculo found, updating existing vehiculo with Matricula: %s", vehiculo.Matricula)
		var existingVehiculo models.Vehiculo
		if err := json.NewDecoder(resp.Body).Decode(&existingVehiculo); err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to decode existing vehiculo response: %s", string(bodyBytes))
		}
		vehiculo.ID = existingVehiculo.ID

		log.Printf("Sending PUT request to update vehiculo with Matricula: %s", vehiculo.Matricula)
		req, err = http.NewRequest("PUT", "http://localhost:8080/vehiculos/"+vehiculo.ID.String(), bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("update vehiculo request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}
		log.Printf("Vehiculo updated successfully with Matricula: %s", vehiculo.Matricula)
	} else {
		log.Printf("Vehiculo already exists with Matricula: %s", vehiculo.Matricula)
	}

	return nil
}

// WebhookHandler maneja las solicitudes entrantes del webhook
func WebhookHandler(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON", "details": err.Error()})
		return
	}

	// Verificar que el payload contenga el campo "mensaje"
	mensaje, ok := payload["mensaje"].(string)
	if !ok {
		log.Printf("Invalid or missing 'mensaje' field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'mensaje' field"})
		return
	}

	// Procesar el mensaje en párrafos
	paragraphs := strings.Split(mensaje, "\n\n")
	log.Printf("Paragraphs extracted from message: %v", paragraphs)

	if len(paragraphs) < 3 {
		log.Printf("Invalid message format: not enough paragraphs")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format", "details": "Message does not contain enough paragraphs"})
		return
	}

	firstParagraph := paragraphs[2]  // El tercer párrafo contiene la información de llegada
	secondParagraph := paragraphs[0] // El primer párrafo contiene el encabezado

	log.Printf("First paragraph: %s", firstParagraph)
	log.Printf("Second paragraph: %s", secondParagraph)

	// Usar una expresión regular para encontrar la fecha en el formato esperado
	re := regexp.MustCompile(`\d{4}:\d{2}[A-Z]{3}\d{2}`)
	dateString := re.FindString(firstParagraph)
	if dateString == "" {
		log.Printf("Date string not found in message")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date string not found in message"})
		return
	}
	log.Printf("Extracted date string: %s", dateString)

	// Crear una cadena de fecha en el formato RFC822
	day := dateString[:2]
	hour := dateString[2:4]
	minute := dateString[5:7]
	month := dateString[7:10]
	year := "20" + dateString[10:] // Asumimos que el año es 20XX

	dateFormatted := fmt.Sprintf("%s %s %s %s:%s UTC", day, month, year, hour, minute)
	log.Printf("Formatted date string: %s", dateFormatted)

	// Convertir la fecha a time.Time
	date, err := time.Parse("02 Jan 2006 15:04 MST", dateFormatted)
	if err != nil {
		log.Printf("Invalid date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "details": err.Error()})
		return
	}
	log.Printf("Parsed date: %v", date)

	// Extraer el lugar
	placeStart := strings.Index(firstParagraph, "por el Aeropuerto")
	if placeStart == -1 {
		log.Printf("'por el Aeropuerto' not found in message")
		c.JSON(http.StatusBadRequest, gin.H{"error": "'por el Aeropuerto' not found in message"})
		return
	}
	placeStart += len("por el Aeropuerto")

	placeEnd := strings.Index(firstParagraph[placeStart:], " de la aeronave")
	if placeEnd == -1 {
		log.Printf("' de la aeronave' not found in message")
		c.JSON(http.StatusBadRequest, gin.H{"error": "' de la aeronave' not found in message"})
		return
	}
	placeEnd += placeStart

	place := strings.TrimSpace(firstParagraph[placeStart:placeEnd])
	log.Printf("Extracted place: %s", place)

	// Crear el objeto IIO
	iio := models.IIO{
		Descripcion: mensaje,
		Fecha:       date,
		Lugar:       place,
		Area:        "ciberseguridad",
		Nombre:      secondParagraph,
	}
	log.Printf("Constructed IIO object: %+v", iio)

	// Realizar solicitud de inicio de sesión para obtener el token JWT
	loginURL := "http://localhost:8080/login"
	loginData := LoginRequest{
		Correo:   "string",
		Password: "string",
	}

	token, err := getJWTToken(loginURL, loginData)
	if err != nil {
		log.Printf("Failed to get JWT token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get JWT token", "details": err.Error()})
		return
	}

	// Añadir el token JWT al encabezado para crear IIO
	iioJsonData, err := json.Marshal(iio)
	if err != nil {
		log.Printf("Failed to marshal IIO data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal IIO data", "details": err.Error()})
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/iios", bytes.NewBuffer(iioJsonData))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request", "details": err.Error()})
		return
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		c.JSON(resp.StatusCode, gin.H{"error": "Request failed", "details": string(bodyBytes)})
		return
	}

	var createdIIO models.IIO
	if err := json.NewDecoder(resp.Body).Decode(&createdIIO); err != nil {
		log.Printf("Failed to decode response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response", "details": err.Error()})
		return
	}

	// Extraer personas del mensaje y procesarlas
	personasRe := regexp.MustCompile(`[A-Za-z]+\s+[A-Za-z]+,\s*V-\d+\.\d+\.\d+`)
	personasMatches := personasRe.FindAllString(mensaje, -1)

	for _, match := range personasMatches {
		parts := strings.Split(match, ",")
		nombre := strings.TrimSpace(parts[0])
		cedula := strings.TrimSpace(parts[1])

		// Limpiar la cédula
		cedula = strings.ReplaceAll(strings.ReplaceAll(cedula, "V-", ""), ".", "")
		log.Printf("Processed person: %s, %s", nombre, cedula)

		persona := models.Persona{
			Nombre:    nombre,
			Cedula:    cedula,
			UserID:    createdIIO.UserID,
			IIOs:      []models.IIO{createdIIO},
			Area:      "ciberseguridad",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		log.Printf("Sending request to create or update persona with Cedula: %s", persona.Cedula)
		if err := createOrUpdatePersona(persona, token); err != nil {
			log.Printf("Error creating or updating persona with Cedula: %s, Error: %v", persona.Cedula, err)
			continue // No detiene el proceso, simplemente pasa al siguiente registro
		}
		log.Printf("Processed persona with Cedula: %s", persona.Cedula)
	}

	// Extraer matrícula del mensaje y procesarla
	matriculaRe := regexp.MustCompile(`matricula\s+([A-Z0-9]+),`)
	matriculaMatch := matriculaRe.FindStringSubmatch(mensaje)

	if len(matriculaMatch) > 1 {
		matricula := matriculaMatch[1]
		log.Printf("Extracted matricula: %s", matricula)

		vehiculo := models.Vehiculo{
			Matricula: matricula,
			UserID:    createdIIO.UserID,
			IIOs:      []models.IIO{createdIIO},
			Area:      "ciberseguridad",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		log.Printf("Sending request to create or update vehiculo with Matricula: %s", vehiculo.Matricula)
		if err := createOrUpdateVehiculo(vehiculo, token); err != nil {
			log.Printf("Error creating or updating vehiculo with Matricula: %s, Error: %v", vehiculo.Matricula, err)
			// No detiene el proceso, simplemente pasa al siguiente registro
		}
		log.Printf("Processed vehiculo with Matricula: %s", vehiculo.Matricula)
	}

	c.JSON(http.StatusOK, createdIIO)
}
