package models

import (
	"time"

	
)

type CreateUserRequest struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Cedula     string `json:"cedula"`
	Telefono   string `json:"telefono"`
	Password   string `json:"password"`
	Credencial string `json:"credencial"`
	Usuario    string         `json:"u_telegran"`
	Alias        string         `json:"alias"`
	Correo     string `json:"correo"`
	REDI       string         `json:"redi"`
	ADI        string         `json:"adi"`
	Zodi        string         `json:"zodi"`
	Area       string `json:"area"`
	Fecha      time.Time      `json:"fecha_nacimiento"`
	Nivel      string `json:"nivel"`
	Descripcion        string `json:"descripcion"`
	Tie        string    	  `json:"tie"`

}

type LoginRequest struct {
	Correo   string `json:"correo" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}