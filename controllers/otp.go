package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	configs "github.com/oficialrivas/sgi/config"
	"github.com/oficialrivas/sgi/models"
	"github.com/pquerna/otp/totp"
)

// OTPResponse representa la estructura de la respuesta OTP
type OTPResponse struct {
	OTPURL string `json:"otp_url"`
}

// ErrorResponse representa la estructura de las respuestas de error
type ErrorResponse struct {
	Error string `json:"error"`
}

// SetupOTP genera un secreto OTP y un URL para configurar Google Authenticator
// @Summary Genera un secreto OTP y un URL para configurar Google Authenticator
// @Description Genera un secreto OTP y proporciona un URL en formato QR para configurar Google Authenticator
// @Tags OTP
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} OTPResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id}/otp-setup [get]
// @Security BearerAuth
func SetupOTP(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := configs.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SIGEIN",
		AccountName: user.Correo,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate OTP key"})
		return
	}

	user.OTPSecret = key.Secret()
	if err := configs.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to save OTP secret"})
		return
	}

	c.JSON(http.StatusOK, OTPResponse{OTPURL: key.URL()})
}
