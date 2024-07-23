package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	// Optional Refresh Token
	RefreshToken string `json:"refreshToken"`
	ID           string `json:"id"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}


// OTPResponse representa la estructura de la respuesta OTP
type OTPResponse struct {
	OTPURL string `json:"otp_url"`
}

