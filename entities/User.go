package entities

import "time"

type User struct {
	Base
	Username    string       `gorm:"not null"`
	Password    string       `gorm:"not null"`
	UserSymbols []UserSymbol `gorm:"foreignKey:UserID"`
	LastSeen    time.Time
}
