package entities

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UserSymbol struct {
	UserID    uuid.UUID `json:"userId" gorm:"column:user_id;primaryKey"`
	SymbolID  uuid.UUID `json:"symbolId" gorm:"column:symbol_id;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
	Amount    float64   `json:"amount"`
	Symbol    *Symbol   `json:"symbol" gorm:"foreignKey:ID;references:SymbolID"`
}
