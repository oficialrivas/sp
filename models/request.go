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

// UpdateTelegramRequest representa la estructura del cuerpo de la solicitud para actualizar el campo u_telegram
type UpdateTelegramRequest struct {
    Usuario string `json:"u_telegram" binding:"required"`
}


type IIORequestParams struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	REDI       string `json:"redi,omitempty"`
	Tie   string `json:"tie,omitempty"`
	Modalidad  string `json:"modalidad,omitempty"`
}

type IIORequestParams2 struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Modalidad  string `json:"modalidad,omitempty"`
	Valor      *bool  `json:"valor,omitempty"`  // Campo booleano para filtrar
}