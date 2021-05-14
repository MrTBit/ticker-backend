package entities

type Symbol struct {
	Base
	Symbol      string `gorm:"size:255;not null"`
	Description string
	Price       float64
	LastPrice   float64
	Active      bool `gorm:"not null"`
}
