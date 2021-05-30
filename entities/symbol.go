package entities

type Symbol struct {
	Base
	Symbol      string  `json:"symbol" gorm:"size:255;not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	LastPrice   float64 `json:"lastPrice"`
	Active      bool    `json:"active" gorm:"not null"`
}
