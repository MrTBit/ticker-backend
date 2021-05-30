package entities

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null" json:"id"`
	CreatedAt time.Time `json:"createdAt,string" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt,string" gorm:"not null"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(_ *gorm.DB) (err error) {
	base.ID = uuid.NewV4()
	return
}
