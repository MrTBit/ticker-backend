package models

import "time"

type User struct {
	Base
	Username string   `gorm:"not null"`
	Password string   `gorm:"not null"`
	Symbol   []Symbol `gorm:"many2many:user_symbols;"`
	LastSeen time.Time
}
