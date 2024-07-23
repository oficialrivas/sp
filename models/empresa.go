package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Empresa struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Nombre          string    `json:"nombre"`
	Direccion       string    `json:"direccion"`
	Actividad       string    `json:"actividad_economica"`
	RIF             string    `gorm:"unique" json:"rif"`
	Area            string    `json:"area"`
	Representante   Persona   `gorm:"foreignKey:RepresentanteID" json:"representante"`
	RepresentanteID uuid.UUID `json:"representante_id"`
	Socios          []Persona `gorm:"many2many:empresa_socios;" json:"socios"`
	Empleados       []Persona `gorm:"many2many:empresa_empleados;" json:"empleados"`
	UserID          uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Empresa) TableName() string {
	return "empresa"
}

func (empresa *Empresa) BeforeCreate(tx *gorm.DB) (err error) {
	empresa.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	empresa.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	empresa.UpdatedAt = empresa.CreatedAt
	return nil
}
