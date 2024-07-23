package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Visa struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Valoracion      string    `json:"valoracion"`
	Pais            string    `json:"pais"`
	Tipo            string    `json:"tipo"`
	Area            string    `json:"area"`
	Codigo          string    `gorm:"unique" json:"codigo"`
	Representante   Persona   `gorm:"foreignKey:RepresentanteID" json:"representante"`
	RepresentanteID uuid.UUID `json:"representante_id"`
	UserID          uuid.UUID `gorm:"type:uuid;column:user_id"`
	Aprobada        bool      `json:"aprobada"`
}

func (Visa) TableName() string {
	return "visa"
}

func (visa *Visa) BeforeCreate(tx *gorm.DB) (err error) {
	visa.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	visa.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	visa.UpdatedAt = visa.CreatedAt
	return nil
}
