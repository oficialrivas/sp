package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Modalidad struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Nombre      string    `json:"nombre"`
	Area            string    `json:"area"`
	UserID          uuid.UUID `gorm:"type:uuid;column:user_id"`
	
}

func (Modalidad) TableName() string {
	return "modalidad"
}

func (modalidad *Modalidad) BeforeCreate(tx *gorm.DB) (err error) {
	modalidad.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	modalidad.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	modalidad.UpdatedAt = modalidad.CreatedAt
	return nil
}
