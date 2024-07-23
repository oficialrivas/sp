package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Caso struct {
	ID          uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Nombre      string      `gorm:"unique" json:"nombre"`
	Tipo        string      `json:"tipo"`
	Codigo      string      `gorm:"unique" json:"codigo"`
	Area        string      `json:"area"`
	Modalidad         string    `json:"modalidad"`
	Tie         string    `json:"tie"`
	Vdirector         int    `json:"vdirector"`
	Vanalista         int    `json:"vanalista"`
	Vcoordinador         int    `json:"vcoordinador"`
	Relacion    []Persona   `gorm:"many2many:relacion_caso;" json:"relacion"`
	Vehiculos   []Vehiculo  `gorm:"many2many:caso_vehiculos;" json:"vehiculos"`
	Empresas    []Empresa   `gorm:"many2many:caso_empresas;" json:"empresas"`
	Direcciones []Direccion `gorm:"many2many:casp_direcciones;" json:"direcciones"`
	IIOs        []IIO       `gorm:"many2many:caso_iio;" json:"iio"`
	Documentos  []Documento `gorm:"many2many:caso_documento;" json:"documento"`
	UserID       uuid.UUID `gorm:"type:uuid;column:user_id"`
	Users       []User      `gorm:"many2many:caso_users;" json:"users"`
}

func (Caso) TableName() string {
	return "caso"
}

func (caso *Caso) BeforeCreate(tx *gorm.DB) (err error) {
	caso.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	caso.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	caso.UpdatedAt = caso.CreatedAt
	return nil
}
