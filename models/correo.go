package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Correo struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Tipo            string    `json:"tipo"`
	Area            string    `json:"area"`
	Direccion       string    `gorm:"unique" json:"direccion"`
	DuenoID         uuid.UUID `json:"dueno_id"`
	Dueno           Persona   `gorm:"foreignKey:DuenoID" json:"dueno"`
	UserID          uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Correo) TableName() string {
	return "correo"
}

func (correo *Correo) BeforeCreate(tx *gorm.DB) (err error) {
	correo.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	correo.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	correo.UpdatedAt = correo.CreatedAt
	return nil
}
