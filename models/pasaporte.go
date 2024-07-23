package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pasaporte struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Numero          string    `json:"numero"`
	Foto            string    `json:"foto"`
	Pais            string    `json:"pais"`
	Tipo            string    `json:"tipo"`
	Area            string    `json:"area"`
	Codigo          string    `gorm:"unique" json:"codigo"`
	Representante   Persona   `gorm:"foreignKey:RepresentanteID" json:"representante"`
	RepresentanteID uuid.UUID `json:"representante_id"`
	UserID          uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Pasaporte) TableName() string {
	return "pasaporte"
}

func (pasaporte *Pasaporte) BeforeCreate(tx *gorm.DB) (err error) {
	pasaporte.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	pasaporte.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	pasaporte.UpdatedAt = pasaporte.CreatedAt
	return nil
}
