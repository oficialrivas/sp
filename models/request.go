package models

type CreateUserRequest struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Cedula     string `json:"cedula"`
	Telefono   string `json:"telefono"`
	Password   string `json:"password"`
	Credencial string `json:"credencial"`
	Correo     string `json:"correo"`
	Area       string `json:"area"`
	Nivel      string `json:"nivel"`
	Tie        string    	  `json:"tie"`
	REDI       string   	  `json:"redi"`
}

type LoginRequest struct {
	Correo   string `json:"correo" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}