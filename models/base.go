package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	CreatedAt time.Time `json:"created_at,string" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at,string" gorm:"not null"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(_ *gorm.DB) (err error) {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		return err
	}
	base.ID = uuidGen
	return
}
