package models

type FHRecvMessage struct {
	Data        []FHRecvMessageData `json:"data"`
	MessageType string              `json:"type"`
}

type FHRecvMessageData struct {
	Price  float64 `json:"p"`
	Symbol string  `json:"s"`
	Time   int64   `json:"t"`
	Volume float64 `json:"v"`
}
