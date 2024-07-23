package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Redes struct {
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

func (Redes) TableName() string {
	return "redes"
}

func (redes *Redes) BeforeCreate(tx *gorm.DB) (err error) {
	redes.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	redes.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	redes.UpdatedAt = redes.CreatedAt
	return nil
}
