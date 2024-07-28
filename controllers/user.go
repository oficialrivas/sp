package controllers

import (
	"net/http"
	"time"
	"log"

	"fmt"
	"strings"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/utils"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"github.com/xuri/excelize/v2"
)



// CreateUser crea un nuevo usuario
// @Summary Crea un nuevo usuario
// @Description Crea un nuevo usuario con los datos proporcionados
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Datos del usuario"
// @Success 200 
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /signup [post]
func CreateUser(c *gin.Context) {
	var request models.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	user := models.User{
		Nombre:      request.Nombre,
		Apellido:    request.Apellido,
		Cedula:      request.Cedula,
		Telefono:    request.Telefono,
		Usuario:     request.Usuario,
		Hash:        string(hashedPassword),
		Credencial:  request.Credencial,
		Correo:      request.Correo,
		Area:        request.Area,
		Alias:       request.Alias,
		Fecha:       request.Fecha,
		Descripcion: request.Descripcion,
		Nivel:       request.Nivel,
		Tie:         request.Tie,
		REDI:        request.REDI,
		Zodi:        request.Zodi,
		ADI:         request.ADI,
	}

	if err := configs.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario creado con éxito"})
}

// Login autentica a un usuario y genera un JWT
// @Summary Autentica a un usuario
// @Description Autentica a un usuario con correo y contraseña, y genera un token JWT
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Credenciales de inicio de sesión"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	if err := configs.DB.Where("correo = ?", request.Correo).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Correo o contraseña incorrectos"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Correo o contraseña incorrectos"})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID.String(), user.Nivel, user.Area)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo generar los tokens"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken, ID: user.ID.String()})
}


// GenerateToken genera un JWT válido a partir del número de teléfono o usuario de Telegram
// @Summary Genera un JWT válido
// @Description Genera un JWT válido a partir del número de teléfono o usuario de Telegram
// @Tags users
// @Accept json
// @Produce json
// @Param data body models.GenerateTokenRequest true "Número de teléfono o usuario de Telegram"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /generate-token [post]
func GenerateToken(c *gin.Context) {
	var request models.GenerateTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	if request.Telefono != "" {
		if err := configs.DB.Where("telefono = ?", request.Telefono).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Usuario no encontrado"})
			return
		}
	} else if request.Usuario != "" {
		if err := configs.DB.Where("usuario = ?", request.Usuario).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Usuario no encontrado"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Se requiere número de teléfono o usuario de Telegram"})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID.String(), user.Nivel, user.Area)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo generar los tokens"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken, ID: user.ID.String()})
}


// RefreshToken renueva el accessToken usando el refreshToken
// @Summary Renueva el accessToken
// @Description Renueva el accessToken usando el refreshToken proporcionado
// @Tags users
// @Accept json
// @Produce json
// @Param tokens body models.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /refresh-token [post]
func RefreshToken(c *gin.Context) {
	var request models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	claims, err := utils.ValidateJWT(request.RefreshToken, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Token de actualización inválido"})
		return
	}

	accessToken, _, err := utils.GenerateTokens(claims.UserID, claims.Role, claims.Area)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{AccessToken: accessToken, ID: claims.UserID})
}

// GetUser obtiene un usuario por su ID
// @Summary Obtiene un usuario por su ID
// @Description Obtiene los datos de un usuario por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del usuario"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
// @Security BearerAuth
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := configs.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers obtiene todos los usuarios
// @Summary Obtiene todos los usuarios
// @Description Obtiene una lista de todos los usuarios
// @Tags users
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
// @Security BearerAuth
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := configs.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser actualiza un usuario por su ID
// @Summary Actualiza un usuario por su ID
// @Description Actualiza los datos de un usuario por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del usuario"
// @Param user body models.CreateUserRequest true "Datos actualizados del usuario"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put]
// @Security BearerAuth
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var request models.CreateUserRequest
	var user models.User
	if err := configs.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user.Nombre = request.Nombre
	user.Apellido = request.Apellido
	user.Cedula = request.Cedula
	user.Telefono = request.Telefono
	user.Credencial = request.Credencial
	user.Correo = request.Correo
	user.Area = request.Area
	user.Nivel = request.Nivel

	if err := configs.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser borra un usuario por su ID
// @Summary Borra un usuario por su ID
// @Description Borra un usuario de la base de datos por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del usuario"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
// @Security BearerAuth
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := configs.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		return
	}

	if err := configs.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User deleted"})
}


// UpdatePasswordRequest estructura para la solicitud de actualización de contraseña
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	OTP         string `json:"otp"`
}

// SuccessResponse representa la estructura de una respuesta exitosa
type SuccessResponse struct {
	Message string `json:"message"`
}

// UpdatePassword actualiza la contraseña de un usuario, verificando OTP para administradores
// @Summary Actualiza la contraseña de un usuario
// @Description Actualiza la contraseña de un usuario verificando el OTP para administradores
// @Tags OTP
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param body body UpdatePasswordRequest true "Datos de la solicitud para actualizar la contraseña"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/password [put]
// @Security BearerAuth
func UpdatePassword(c *gin.Context) {
	userID := c.Param("id")
	var request UpdatePasswordRequest
	var user models.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := configs.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(request.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Incorrect password"})
		return
	}

	if user.Nivel == "admin" {
		if !totp.Validate(request.OTP, user.OTPSecret) {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid OTP"})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
		return
	}

	user.Hash = string(hashedPassword)
	user.UpdatedAt = time.Now().UTC()

	if err := configs.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Password updated successfully"})
}

// GetMensajesByUserID obtiene los mensajes asociados a un ID de usuario
// @Summary Obtiene los mensajes asociados a un ID de usuario
// @Description Obtiene todos los mensajes asociados al ID de un usuario proporcionado
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del usuario"
// @Success 200 {array} models.Mensaje
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id}/messages [get]
// @Security ApiKeyAuth
func GetMensajesByUserID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User ID is required"})
		return
	}

	var mensajes []models.Mensaje
	if err := configs.DB.Where("user_id = ?", userID).Find(&mensajes).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Mensajes no encontrados"})
		return
	}

	c.JSON(http.StatusOK, mensajes)
}


// GetUsersByNivel obtiene los usuarios por nivel
// @Summary Obtiene los usuarios por nivel
// @Description Obtiene una lista de usuarios que coinciden con el nivel proporcionado
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param nivel query string true "Nivel del usuario (admin, superuser, analyst, user)"
// @Success 200 {array} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/nivel [get]
// @Security BearerAuth
func GetUsersByNivel(c *gin.Context) {
	nivel := c.Query("nivel")
	if nivel == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Nivel es requerido"})
		return
	}

	validNiveles := map[string]bool{
		"admin":     true,
		"superuser": true,
		"analyst":   true,
		"user":      true,
	}

	if !validNiveles[nivel] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Nivel no es válido"})
		return
	}

	var users []models.User
	if err := configs.DB.Where("nivel = ?", nivel).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUsersWithUnprocessedMessages obtiene todos los usuarios que tienen mensajes no procesados
// @Summary Obtiene todos los usuarios con mensajes no procesados
// @Description Obtiene una lista de todos los usuarios que tienen mensajes no procesados
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users-with-unprocessed-messages [get]
// @Security BearerAuth
func GetUsersWithUnprocessedMessages(c *gin.Context) {
	// Obtener todos los mensajes no procesados
	var mensajes []models.Mensaje
	if err := configs.DB.Where("procesado = ?", false).Find(&mensajes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Obtener los IDs de los usuarios a partir de los mensajes no procesados
	userIDsMap := make(map[uuid.UUID]bool)
	for _, mensaje := range mensajes {
		userIDsMap[mensaje.UserID] = true
	}

	// Convertir el mapa a una lista de IDs
	var userIDs []uuid.UUID
	for userID := range userIDsMap {
		userIDs = append(userIDs, userID)
	}

	// Obtener los detalles de los usuarios usando los IDs
	var users []models.User
	if err := configs.DB.Where("id IN (?)", userIDs).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}


// GetUsersWithUnprocessedMessagesByREDI obtiene todos los usuarios de una REDI específica que tienen mensajes no procesados
// @Summary Obtiene todos los usuarios de una REDI específica con mensajes no procesados
// @Description Obtiene una lista de todos los usuarios de una REDI específica que tienen mensajes no procesados
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param redi path string true "REDI"
// @Success 200 {array} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users-with-unprocessed-messages-by-redi/{redi} [get]
// @Security BearerAuth
func GetUsersWithUnprocessedMessagesByREDI(c *gin.Context) {
	// Obtener el valor de la REDI desde los parámetros de la URL
	redi := c.Param("redi")
	if redi == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "REDI is required"})
		return
	}

	// Obtener todos los usuarios de la REDI especificada
	var users []models.User
	if err := configs.DB.Where("redi = ?", redi).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Filtrar usuarios que tienen mensajes no procesados
	var usersWithUnprocessedMessages []models.User
	for _, user := range users {
		var count int64
		if err := configs.DB.Model(&models.Mensaje{}).Where("user_id = ? AND procesado = ?", user.ID, false).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		if count > 0 {
			usersWithUnprocessedMessages = append(usersWithUnprocessedMessages, user)
		}
	}

	c.JSON(http.StatusOK, usersWithUnprocessedMessages)
}


// GetUsersWithUnprocessedMessagesByREDI obtiene todos los usuarios de una REDI específica con mensajes no procesados
// @Summary Obtiene todos los usuarios de una REDI específica con mensajes no procesados
// @Description Obtiene una lista de todos los usuarios de una REDI específica con mensajes no procesados
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param redi path string true "REDI"
// @Success 200 {array} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users-with-unprocessed-messages-by-redi-and-nivel/{redi} [get]
// @Security BearerAuth
func GetUsersWithUnprocessedMessagesByREDIAndNivel(c *gin.Context) {
	// Obtener el valor de la REDI desde los parámetros de la URL
	redi := c.Param("redi")
	if redi == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "REDI is required"})
		return
	}

	log.Printf("REDI provided: %s", redi)

	// Obtener todos los usuarios de la REDI especificada y con nivel "user"
	var users []models.User
	if err := configs.DB.Where("redi = ? AND nivel = ?", redi, "user").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Verificar si se encontraron usuarios
	if len(users) == 0 {
		log.Printf("No users found for REDI: %s", redi)
		c.JSON(http.StatusOK, gin.H{"message": "No se encontraron usuarios con la REDI especificada y el nivel de usuario"})
		return
	}

	log.Printf("Users found: %v", users)

	// Filtrar usuarios que tienen mensajes no procesados
	var usersWithUnprocessedMessages []models.User
	for _, user := range users {
		var count int64
		if err := configs.DB.Model(&models.Mensaje{}).Where("user_id = ? AND procesado = ?", user.ID, false).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		log.Printf("User ID: %s, Unprocessed Messages Count: %d", user.ID, count)
		if count > 0 {
			usersWithUnprocessedMessages = append(usersWithUnprocessedMessages, user)
		}
	}

	// Verificar si se encontraron usuarios con mensajes no procesados
	if len(usersWithUnprocessedMessages) == 0 {
		log.Printf("No unprocessed messages found for users in REDI: %s", redi)
		c.JSON(http.StatusOK, gin.H{"message": "No se encontraron mensajes no procesados para los usuarios de la REDI especificada"})
		return
	}

	log.Printf("Users with unprocessed messages: %v", usersWithUnprocessedMessages)

	c.JSON(http.StatusOK, usersWithUnprocessedMessages)
}


// GetUsersWithUnprocessedMessages obtiene todos los usuarios con nivel "user" que tienen mensajes no procesados
// @Summary Obtiene todos los usuarios con nivel "user" que tienen mensajes no procesados
// @Description Obtiene una lista de todos los usuarios con nivel "user" que tienen mensajes con el campo procesado en false
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users-unprocessed-messages-user [get]
// @Security BearerAuth
func GetMensajeUser(c *gin.Context) {
	var users []models.User

	// Subconsulta para obtener IDs de usuarios con mensajes no procesados
	subQuery := configs.DB.Model(&models.Mensaje{}).Select("user_id").Where("procesado = ?", false).Group("user_id")

	// Consulta principal para obtener detalles de los usuarios con nivel "user"
	if err := configs.DB.Where("id IN (?) AND nivel = ?", subQuery, "user").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Verificar si se encontraron usuarios
	if len(users) == 0 {
		log.Println("No se encontraron usuarios con mensajes no procesados y nivel 'user'")
		c.JSON(http.StatusOK, gin.H{"message": "No se encontraron usuarios con mensajes no procesados y nivel 'user'"})
		return
	}

	c.JSON(http.StatusOK, users)
}



// UpdateUserTelegram actualiza el u_telegram de un usuario por su número de teléfono
// @Summary Actualiza el u_telegram de un usuario por su número de teléfono
// @Description Actualiza el campo u_telegram de un usuario utilizando su número de teléfono
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param telefono path string true "Número de teléfono del usuario"
// @Param body body models.UpdateTelegramRequest true "Nuevo u_telegram del usuario"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/telefono/{telefono} [put]
// @Security BearerAuth
func UpdateUserTelegram(c *gin.Context) {
	telefono := c.Param("telefono")
	var request models.UpdateTelegramRequest
	var user models.User

	// Buscar el usuario por número de teléfono
	if err := configs.DB.First(&user, "telefono = ?", telefono).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		return
	}

	// Enlazar el JSON del request a la estructura UpdateTelegramRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Actualizar el campo u_telegram del usuario
	user.Usuario = request.Usuario

	// Guardar los cambios en la base de datos
	if err := configs.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}


// UploadUsersExcel carga usuarios desde un archivo Excel
// @Summary Carga usuarios desde un archivo Excel
// @Description Carga usuarios desde un archivo Excel
// @Tags users
// @Accept multipart/form-data
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Archivo Excel"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /upload_users [post]
func UploadUsersExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "No se proporcionó ningún archivo"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo abrir el archivo"})
		return
	}
	defer f.Close()

	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo leer el archivo Excel"})
		return
	}

	// Verificar que la hoja existe
	sheetName := "Sheet1" // Cambia esto si el nombre de la hoja es diferente
	sheets := excelFile.GetSheetMap()
	found := false
	for _, name := range sheets {
		if strings.EqualFold(name, sheetName) {
			found = true
			sheetName = name
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("La hoja %s no existe en el archivo Excel", sheetName)})
		return
	}

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo obtener las filas del archivo Excel"})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "El archivo Excel está vacío"})
		return
	}

	// Mapa para las posiciones de las columnas
	columnMap := make(map[string]int)
	for i, colName := range rows[0] {
		columnMap[strings.ToLower(colName)] = i
	}

	var processedUsers []models.User

	for _, row := range rows[1:] {
		var request models.CreateUserRequest

		if idx, exists := columnMap["cedula"]; exists && len(row) > idx {
			request.Cedula = strings.ReplaceAll(row[idx], ".", "")
		}
		if idx, exists := columnMap["nombre"]; exists && len(row) > idx {
			request.Nombre = row[idx]
		}
		if idx, exists := columnMap["apellido"]; exists && len(row) > idx {
			request.Apellido = row[idx]
		}
		if idx, exists := columnMap["telefono"]; exists && len(row) > idx {
			request.Telefono = row[idx]
		}
		if idx, exists := columnMap["usuario"]; exists && len(row) > idx {
			request.Usuario = row[idx]
		}
		if idx, exists := columnMap["alias"]; exists && len(row) > idx {
			request.Alias = row[idx]
		}
		if idx, exists := columnMap["descripcion"]; exists && len(row) > idx {
			request.Descripcion = row[idx]
		}
		if idx, exists := columnMap["area"]; exists && len(row) > idx {
			request.Area = row[idx]
		}
		if idx, exists := columnMap["tie"]; exists && len(row) > idx {
			request.Tie = row[idx]
		}
		if idx, exists := columnMap["redi"]; exists && len(row) > idx {
			request.REDI = row[idx]
		}
		if idx, exists := columnMap["zodi"]; exists && len(row) > idx {
			request.Zodi = row[idx]
		}
		if idx, exists := columnMap["adi"]; exists && len(row) > idx {
			request.ADI = row[idx]
		}
		if idx, exists := columnMap["credencial"]; exists && len(row) > idx {
			request.Credencial = row[idx]
		}
		if idx, exists := columnMap["correo"]; exists && len(row) > idx {
			request.Correo = row[idx]
		}
		if idx, exists := columnMap["fecha"]; exists && len(row) > idx {
			request.Fecha, _ = time.Parse("2006-01-02", row[idx])
		}

		// Asumimos una contraseña predeterminada para los usuarios cargados desde Excel
		request.Password = "defaultPassword"
		// El nivel siempre será "user"
		request.Nivel = "user"

		// Buscar usuario existente por cédula
		var user models.User
		if err := configs.DB.Where("cedula = ?", request.Cedula).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Crear nuevo usuario si no existe
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
				user = models.User{
					ID:          uuid.New(),
					Cedula:      request.Cedula,
					Nombre:      request.Nombre,
					Apellido:    request.Apellido,
					Telefono:    request.Telefono,
					Usuario:     request.Usuario,
					Hash:        string(hashedPassword),
					Credencial:  request.Credencial,
					Correo:      request.Correo,
					Area:        request.Area,
					Alias:       request.Alias,
					Fecha:       request.Fecha,
					Descripcion: request.Descripcion,
					Nivel:       request.Nivel, // Se asegura que el nivel sea "user"
					Tie:         request.Tie,
					REDI:        request.REDI,
					Zodi:        request.Zodi,
					ADI:         request.ADI,
				}
				if err := configs.DB.Create(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: fmt.Sprintf("Error al crear el usuario: %v", err)})
					return
				}
				processedUsers = append(processedUsers, user)
			} else {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: fmt.Sprintf("Error al buscar la TIE: %v", err)})
				return
			}
		} else {
			// Actualizar usuario existente
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
			user.Nombre = request.Nombre
			user.Apellido = request.Apellido
			user.Telefono = request.Telefono
			user.Usuario = request.Usuario
			user.Hash = string(hashedPassword)
			user.Credencial = request.Credencial
			user.Correo = request.Correo
			user.Area = request.Area
			user.Alias = request.Alias
			user.Fecha = request.Fecha
			user.Descripcion = request.Descripcion
			user.Nivel = "user" // Se asegura que el nivel sea "user"
			user.Tie = request.Tie
			user.REDI = request.REDI
			user.Zodi = request.Zodi
			user.ADI = request.ADI
			if err := configs.DB.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: fmt.Sprintf("Error al actualizar el usuario: %v", err)})
				return
			}
			processedUsers = append(processedUsers, user)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuarios procesados correctamente", "users": processedUsers})
}


// DeleteUsersByCedula elimina usuarios desde un archivo Excel
// @Summary Elimina usuarios desde un archivo Excel
// @Description Elimina usuarios desde un archivo Excel
// @Tags users
// @Accept multipart/form-data
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Archivo Excel"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /delete_users [post]
func DeleteUsersByCedula(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "No se proporcionó ningún archivo"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo abrir el archivo"})
		return
	}
	defer f.Close()

	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo leer el archivo Excel"})
		return
	}

	// Verificar que la hoja existe
	sheetName := "Sheet1" // Cambia esto si el nombre de la hoja es diferente
	sheets := excelFile.GetSheetMap()
	found := false
	for _, name := range sheets {
		if strings.EqualFold(name, sheetName) {
			found = true
			sheetName = name
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("La hoja %s no existe en el archivo Excel", sheetName)})
		return
	}

	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "No se pudo obtener las filas del archivo Excel"})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "El archivo Excel está vacío"})
		return
	}

	// Mapa para las posiciones de las columnas
	columnMap := make(map[string]int)
	for i, colName := range rows[0] {
		columnMap[strings.ToLower(colName)] = i
	}

	var deletedUsers []string

	for _, row := range rows[1:] {
		var cedula string
		if idx, exists := columnMap["cedula"]; exists && len(row) > idx {
			cedula = strings.ReplaceAll(row[idx], ".", "")
		}

		if cedula == "" {
			continue // Saltar filas sin cédula
		}

		// Buscar usuario existente por cédula
		var user models.User
		if err := configs.DB.Where("cedula = ?", cedula).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue // Si no se encuentra el usuario, continuar con la siguiente fila
			} else {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: fmt.Sprintf("Error al buscar el usuario: %v", err)})
				return
			}
		}

		// Eliminar usuario existente
		if err := configs.DB.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: fmt.Sprintf("Error al eliminar el usuario: %v", err)})
			return
		}

		deletedUsers = append(deletedUsers, cedula)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuarios eliminados correctamente", "deleted_users": deletedUsers})
}