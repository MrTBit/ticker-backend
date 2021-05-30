package entities

import "time"

type User struct {
	Base
	Username    string       `json:"username" gorm:"not null"`
	Password    string       `json:"password" gorm:"not null"`
	UserSymbols []UserSymbol `json:"userSymbols" gorm:"foreignKey:UserID"`
	LastSeen    time.Time    `json:"lastSeen"`
}
