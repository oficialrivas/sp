package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Mensaje struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	REDI         string    `json:"redi"`
	ZODI         string    `json:"zodi"`
	ADI         string    `json:"adi"`
	Modalidad    string    `json:"modalidad"`
	Tipo          string    `json:"tipo"`
	Tie          string    `json:"tie"`
	Fecha        time.Time `json:"fecha"`
	Lugar        string    `json:"lugar"`
	Parroquia    string    `json:"parroquia"`
	Descripcion  string    `json:"descripcion"`
	Urbanizacion string    `json:"urbanizacion"`
	Nombre       string    `json:"nombre"`
	Area         string    `json:"area"`
	Procesado   bool      `json:"procesado"`
	ImagenURL    string    `json:"imagen_url"`
	Relacion     []Persona `gorm:"many2many:relacion_persona;" json:"relacion"`
	Nivel        string    `json:"nivel"`
	UserID       uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Mensaje) TableName() string {
	return "mensaje"
}

func (mensaje *Mensaje) BeforeCreate(tx *gorm.DB) (err error) {
	mensaje.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	mensaje.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	mensaje.UpdatedAt = mensaje.CreatedAt
	return nil
}
