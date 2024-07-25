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
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/models"
)

func createOrUpdatePersona(persona models.Persona, token string) error {
	jsonData, err := json.Marshal(persona)
	if err != nil {
		return err
	}

	client := &http.Client{}

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
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("create persona request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}
		log.Printf("Persona created successfully with Cedula: %s", persona.Cedula)
	} else if resp.StatusCode == http.StatusOK {
		var existingPersona models.Persona
		if err := json.NewDecoder(resp.Body).Decode(&existingPersona); err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to decode existing persona response: %s", string(bodyBytes))
		}
		persona.ID = existingPersona.ID

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
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("update persona request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}
		log.Printf("Persona updated successfully with Cedula: %s", persona.Cedula)
	} else {
		log.Printf("Unexpected response status: %d", resp.StatusCode)
	}

	return nil
}

func createOrUpdateVehiculo(vehiculo models.Vehiculo, token string) error {
	jsonData, err := json.Marshal(vehiculo)
	if err != nil {
		return err
	}

	client := &http.Client{}

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
		var existingVehiculo models.Vehiculo
		if err := json.NewDecoder(resp.Body).Decode(&existingVehiculo); err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to decode existing vehiculo response: %s", string(bodyBytes))
		}
		vehiculo.ID = existingVehiculo.ID

		log.Printf("Vehiculo found, updating existing vehiculo with Matricula: %s", vehiculo.Matricula)
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
		log.Printf("Unexpected response status: %d", resp.StatusCode)
	}

	return nil
}

// WebhookHandler maneja las solicitudes entrantes del webhook
func WebhookHandler(c *gin.Context) {
	var payload struct {
		Mensaje string `json:"mensaje"`
		JWT     string `json:"jwt"`
		IIOID   string `json:"iio_id"`
	}
	if err := c.BindJSON(&payload); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON", "details": err.Error()})
		return
	}

	// Verificar que el payload contenga el campo "mensaje"
	mensaje := payload.Mensaje
	log.Printf("Received message: %s", mensaje)

	// Usar una expresión regular para encontrar la fecha en el formato esperado
	re := regexp.MustCompile(`\d{4}:\d{2}[A-Z]{3}\d{2}`)
	dateString := re.FindString(mensaje)
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
	placeStart := strings.Index(mensaje, "por el Aeropuerto")
	if placeStart == -1 {
		log.Printf("'por el Aeropuerto' not found in message")
		c.JSON(http.StatusBadRequest, gin.H{"error": "'por el Aeropuerto' not found in message"})
		return
	}
	placeStart += len("por el Aeropuerto")

	placeEnd := strings.Index(mensaje[placeStart:], " de la aeronave")
	if placeEnd == -1 {
		log.Printf("' de la aeronave' not found in message")
		c.JSON(http.StatusBadRequest, gin.H{"error": "' de la aeronave' not found in message"})
		return
	}
	placeEnd += placeStart

	place := strings.TrimSpace(mensaje[placeStart:placeEnd])
	log.Printf("Extracted place: %s", place)

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
			UserID:    uuid.MustParse(payload.IIOID), // Vincular con la IIO
			Area:      "ciberseguridad",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		log.Printf("Sending request to create or update persona with Cedula: %s", persona.Cedula)
		if err := createOrUpdatePersona(persona, payload.JWT); err != nil {
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
			UserID:    uuid.MustParse(payload.IIOID), // Vincular con la IIO
			Area:      "ciberseguridad",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		log.Printf("Sending request to create or update vehiculo with Matricula: %s", vehiculo.Matricula)
		if err := createOrUpdateVehiculo(vehiculo, payload.JWT); err != nil {
			log.Printf("Error creating or updating vehiculo with Matricula: %s, Error: %v", vehiculo.Matricula, err)
			// No detiene el proceso, simplemente pasa al siguiente registro
		}
		log.Printf("Processed vehiculo with Matricula: %s", vehiculo.Matricula)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entities processed successfully"})
}
