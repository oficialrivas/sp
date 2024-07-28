package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IIO struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	REDI         string      `json:"redi"`
	ZODI         string      `json:"zodi"`
	ADI          string      `json:"adi"`
	Modalidad    []Modalidad `gorm:"many2many:iio_modalidades;" json:"modalidad"`
	Tie          []Tie       `gorm:"many2many:iio_ties;" json:"TIE"`
	Fecha        time.Time   `json:"fecha"`
	Lugar        string      `json:"lugar"`
	Parroquia    string      `json:"parroquia"`
	Descripcion  string      `json:"descripcion"`
	Urbanizacion string      `json:"urbanizacion"`
	Nombre       string      `json:"nombre"`
	Area         string      `json:"area"`
	Procesado    bool        `json:"procesado"`
	ImagenURL    string      `json:"imagen_url"`
	Mensaje      []Mensaje   `gorm:"many2many:relacion_mensaje;" json:"mensajes"`
	Relacion     []Persona   `gorm:"many2many:relacion_persona;" json:"relacion"`
	Nivel        string      `json:"nivel"`
	UserID       uuid.UUID   `gorm:"type:uuid;column:user_id"`
}

func (IIO) TableName() string {
	return "iio"
}

func (iio *IIO) BeforeCreate(tx *gorm.DB) (err error) {
	iio.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339)
	iio.CreatedAt, err = time.Parse(time.RFC3339, now)
	if err != nil {
		return err
	}
	iio.UpdatedAt = iio.CreatedAt

	// Verificar que la Modalidad corresponda a la Tie
	if !ModalidadCorrespondeATie(iio.Modalidad, iio.Tie) {
		return errors.New("la modalidad no corresponde a la tie")
	}

	return nil
}

func ModalidadCorrespondeATie(modalidades []Modalidad, ties []Tie) bool {
	// Mapa para verificar correspondencia
	tieMap := make(map[string]bool)
	for _, tie := range ties {
		tieMap[tie.Nombre] = true
	}

	for _, modalidad := range modalidades {
		if !tieMap[modalidad.Nombre] {
			return false
		}
	}

	return true
}
