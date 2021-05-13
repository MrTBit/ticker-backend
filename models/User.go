package models

import "time"

type User struct {
	Base
	Username    string       `gorm:"not null"`
	Password    string       `gorm:"not null"`
	UserSymbols []UserSymbol `gorm:"many2many:user_symbols"`
	LastSeen    time.Time
}
