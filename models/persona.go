package models

import (
	
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


// Persona represents a person in the database
type Persona struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Nombre       string      `json:"nombre"`
	Apellido     string      `json:"apellido"`
	Nacionalidad []Nacionalidad `gorm:"many2many:persona_nacionalidad;" json:"nacionalidad"`
	EstadoCivil  string      `json:"estado_civil"`
	Cedula       string      `gorm:"unique" json:"cedula"`
	Correo       string      `gorm:"unique" json:"correo"`
	Telefono     string      `json:"telefono"`
	Profesion    string      `json:"profesion"`
	Ideologia    string      `json:"ideologia"`
	Cargo        string      `json:"cargo_actual"`
	Alias        string      `json:"alias"`
	Filiacion    string      `json:"filiacion_politica"`
	Religion     string      `json:"religion"`
	Tipo         string      `json:"tipo_perfil"`
	Interes      string      `json:"informacion_de_interes"`
	Valoraciones string      `json:"valoraciones"`
	Area         string      `json:"area"`
	UserID       uuid.UUID   `gorm:"type:uuid;column:user_id"`
	Vehiculos    []Vehiculo  `gorm:"many2many:persona_vehiculos;" json:"vehiculos"`
	Empresas     []Empresa   `gorm:"many2many:persona_empresas;" json:"empresas"`
	Direcciones  []Direccion `gorm:"many2many:persona_direcciones;" json:"direcciones"`
	IIOs         []IIO       `gorm:"many2many:persona_iio;" json:"iio"`
	Visas        []Visa      `gorm:"many2many:persona_visa;" json:"visa"`
	Pasaportes   []Pasaporte `gorm:"many2many:persona_pasaporte;" json:"pasaporte"`
	Correos      []Correo    `gorm:"many2many:persona_correo;" json:"correos"`
	Redes        []Redes     `gorm:"many2many:persona_redes;" json:"redes"`
	Relacionados []Persona   `gorm:"many2many:persona_relacionados;joinForeignKey:ID;joinReferences:ID" json:"relacionados"`
}

// TableName overrides the default table name for Persona
func (Persona) TableName() string {
	return "persona"
}

// BeforeCreate is a GORM hook that runs before a new record is inserted into the database
func (persona *Persona) BeforeCreate(tx *gorm.DB) (err error) {
	persona.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	persona.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	persona.UpdatedAt = persona.CreatedAt
	return nil
}
