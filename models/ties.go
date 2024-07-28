package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tie struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Nombre     		string    `json:"nombre"`
	Area            string    `json:"area"`
	Modalidad       []Modalidad      `gorm:"many2many:tie_modalidad;" json:"modalidad"`
	UserID          uuid.UUID `gorm:"type:uuid;column:user_id"`
	
}

func (Tie) TableName() string {
	return "tie"
}

func (tie *Tie) BeforeCreate(tx *gorm.DB) (err error) {
	tie.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	tie.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	tie.UpdatedAt = tie.CreatedAt
	return nil
}
