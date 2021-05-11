package models

import "time"

type User struct {
	Base
	Username   string `gorm:"not null"`
	Password   string `gorm:"not null"`
	UserSymbol []UserSymbol
	LastSeen   time.Time
}
