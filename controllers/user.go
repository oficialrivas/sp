package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/oficialrivas/sgi/utils"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser crea un nuevo usuario
// @Summary Crea un nuevo usuario
// @Description Crea un nuevo usuario con los datos proporcionados
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Datos del usuario"
// @Success 200 {object} gin.H{"message": "Usuario creado con éxito"}
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