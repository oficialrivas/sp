package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemporaryAccess struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`
	EntityType string    `gorm:"type:varchar(50);not null" json:"entity_type"`
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (TemporaryAccess) TableName() string {
	return "temporary_access"
}

func (access *TemporaryAccess) BeforeCreate(tx *gorm.DB) (err error) {
	access.ID = uuid.New()
	now := time.Now().UTC()
	access.CreatedAt = now
	access.UpdatedAt = now
	return nil
}

func (access *TemporaryAccess) BeforeUpdate(tx *gorm.DB) (err error) {
	access.UpdatedAt = time.Now().UTC()
	return nil
}
