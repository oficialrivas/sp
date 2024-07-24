package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Nombre     string         `json:"nombre"`
	Apellido   string         `json:"apellido"`
	REDI       string         `json:"redi"`
	ADI        string         `json:"adi"`
	Zodi        string         `json:"zodi"`
	Fecha      time.Time      `json:"fecha_nacimiento"`
	Parroquia  string         `json:"parroquia"`
	Tie        string         `json:"tie"`
	Alias        string         `json:"alias"`
	Descripcion        string `json:"descripcion"`
	Cedula     string         `gorm:"unique" json:"cedula"`
	Telefono   string         `gorm:"unique" json:"telefono"`
	Usuario    string         `gorm:"unique" json:"u_telegran"`
	Hash       string         `json:"hash"`
	Credencial string         `gorm:"unique" json:"credencial"`
	Correo     string         `gorm:"unique" json:"correo"`
	Area       string         `json:"area"`
	Nivel      string         `json:"nivel"`
	OTPSecret  string         `json:"otp_secret"` 
}

type GenerateTokenRequest struct {
	Telefono string `json:"telefono,omitempty"`
	Usuario  string `json:"usuario,omitempty"`
}

var ValidAreas = []string{
	"SEP",
	"CI2",
	"TIC",
	"CI2 ESPECIAL",
	"Despacho",
}

func (User) TableName() string {
	return "user"
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	user.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	user.UpdatedAt = user.CreatedAt
	return nil
}
