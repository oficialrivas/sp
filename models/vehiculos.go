package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehiculo struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Modelo        string    `json:"modelo"`
	Tipo          string    `json:"tipo"`
	Linea         string    `json:"linea"`
	Marca         string    `json:"marca"`
	Color         string    `json:"color"`
	Ano           string    `json:"año"`
	Numero        string    `json:"numero"`
	Area          string    `json:"area"`
	Matricula     string    `gorm:"unique" json:"matricula"`
	Propietario    []Persona   `gorm:"many2many:propietario_caso;" json:"propietario"`
	IIOs        []IIO       `gorm:"many2many:vehiculo_iio;" json:"vehiculo"`
	Usuarios      []Persona `gorm:"many2many:usuario_vehiculos;" json:"usuarios"`
	UserID        uuid.UUID `gorm:"type:uuid;column:user_id"`
}

func (Vehiculo) TableName() string {
	return "vehiculo"
}

func (vehiculo *Vehiculo) BeforeCreate(tx *gorm.DB) (err error) {
	vehiculo.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	vehiculo.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	vehiculo.UpdatedAt = vehiculo.CreatedAt
	return nil
}
