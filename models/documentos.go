package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Documento struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Numero    string    `json:"numero"`
	Documento string    `json:"documento"`
	Nombre    string    `json:"nombre"`
	Tipo      string    `json:"tipo"`
	Area      string    `json:"area"`
	Codigo    string    `gorm:"unique" json:"codigo"`
	Relacion  []Persona `gorm:"many2many:relacion_persona;" json:"relacion"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Documento) TableName() string {
	return "documento"
}

func (documento *Documento) BeforeCreate(tx *gorm.DB) (err error) {
	documento.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	documento.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	documento.UpdatedAt = documento.CreatedAt
	return nil
}
