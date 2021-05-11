package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UserSymbol struct {
	SymbolID  uuid.UUID `gorm:"primaryKey;not null"`
	UserID    uuid.UUID `gorm:"primaryKey;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Amount    float64
	User      User
	Symbol    Symbol
}
