package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Nacionalidad struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Gentilicio            string    `json:"gentilicio"`
	Area            string    `json:"area"`
	UserID          uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Nacionalidad) TableName() string {
	return "nacionalidad"
}

func (nacionalidad *Nacionalidad) BeforeCreate(tx *gorm.DB) (err error) {
	nacionalidad.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	nacionalidad.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	nacionalidad.UpdatedAt = nacionalidad.CreatedAt
	return nil
}
