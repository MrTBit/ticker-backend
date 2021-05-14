package entities

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UserSymbol struct {
	UserID    uuid.UUID `gorm:"column:user_id;primaryKey"`
	SymbolID  uuid.UUID `gorm:"column:symbol_id;primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Amount    float64
	Symbol    *Symbol `gorm:"foreignKey:ID;references:SymbolID"`
}
