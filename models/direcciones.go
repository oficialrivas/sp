package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Direccion struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Nombre       string    `json:"nombre"`
	Lugar        string    `json:"lugar"`
	Urbanizacion string    `json:"urbanizacion"`
	Parroquia    string    `json:"parroquia"`
	Estado       string    `json:"estado"`
	Municipio    string    `json:"municipio"`
	Area         string    `json:"area"`
	Dueno        Persona   `gorm:"foreignKey:DuenoID" json:"dueno"`
	DuenoID      uuid.UUID `json:"dueno_id"`
	Usuarios     []Persona `gorm:"many2many:direccion_usuarios;" json:"usuarios"`
	Empleados    []Persona `gorm:"many2many:direccion_empleados;" json:"empleados"`
	UserID       uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Direccion) TableName() string {
	return "direccion"
}

func (direccion *Direccion) BeforeCreate(tx *gorm.DB) (err error) {
	direccion.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	direccion.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	direccion.UpdatedAt = direccion.CreatedAt
	return nil
}
