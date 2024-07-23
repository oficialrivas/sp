package utils

import (
	"errors"
	"os"
	"time"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var jwtKey []byte
var refreshKey []byte

func init() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
	refreshKey = []byte(os.Getenv("REFRESH_SECRET"))
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Area   string `json:"area"`
	jwt.StandardClaims
}

func GenerateTokens(userID, role, area string) (string, string, error) {
	accessToken, err := generateJWT(userID, role, area, 15*time.Minute, jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateJWT(userID, role, area, 7*24*time.Hour, refreshKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateJWT(userID, role, area string, duration time.Duration, key []byte) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Area:   area,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, isRefreshToken bool) (*Claims, error) {
	claims := &Claims{}
	key := jwtKey
	if isRefreshToken {
		key = refreshKey
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
