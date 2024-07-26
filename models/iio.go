package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IIO struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	REDI         string    `json:"redi"`
	ZODI         string    `json:"zodi"`
	ADI         string    `json:"adi"`
	Modalidad    string    `json:"modalidad"`
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
	Mensaje     []Mensaje `gorm:"many2many:relacion_mensaje;" json:"mensajes"`
	Relacion     []Persona `gorm:"many2many:relacion_persona;" json:"relacion"`
	Nivel        string    `json:"nivel"`
	UserID       uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (IIO) TableName() string {
	return "iio"
}

func (iio *IIO) BeforeCreate(tx *gorm.DB) (err error) {
	iio.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	iio.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	iio.UpdatedAt = iio.CreatedAt
	return nil
}
