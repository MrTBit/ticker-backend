package socket

import (
	"github.com/gorilla/websocket"
	"log"
	url "net/url"
	"ticker-backend/database"
	"ticker-backend/entities"
	"ticker-backend/models"
)

var conn *websocket.Conn

func InitSocket(interrupt <-chan models.SocketInterrupt) {
	wsURL := url.URL{Scheme: "wss", Host: "ws.finnhub.io", RawQuery: "token=c035nmn48v6v2t3i3n00"}

	var err error
	conn, _, err = websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	SubscribeToAllActiveSymbols()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var messageJson models.FHRecvMessage
			err := conn.ReadJSON(&messageJson)
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Println("recv: ", messageJson)

			if messageJson.MessageType == "trade" {
				handleMessageRecv(messageJson)
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case op := <-interrupt:
			if op.InterruptType == "subscribe" {
				subscribe(op.Symbol)
			} else if op.InterruptType == "unsubscribe" {
				unsubscribe(op.Symbol)
			}
		}
	}
}

func handleMessageRecv(message models.FHRecvMessage) {

	var receivedFilteredData []models.FHRecvMessageData

	//get first of all different symbols
	for _, data := range message.Data {
		foundMatch := false
		for _, filteredData := range receivedFilteredData {
			if filteredData.Symbol == data.Symbol {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			receivedFilteredData = append(receivedFilteredData, data)
		}
	}

	var symbolToUpdate entities.Symbol
	//get and update all symbols with new price
	for _, data := range receivedFilteredData {
		symbolToUpdate = entities.Symbol{}
		database.DBConn.Where("symbol = ?", data.Symbol).First(&symbolToUpdate)

		symbolToUpdate.LastPrice = symbolToUpdate.Price
		symbolToUpdate.Price = data.Price

		database.DBConn.Save(&symbolToUpdate)
	}

}

func subscribe(symbol string) {
	err := conn.WriteJSON(models.FHSendMessage{MessageType: "subscribe", Symbol: symbol})
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func unsubscribe(symbol string) {
	err := conn.WriteJSON(models.FHSendMessage{MessageType: "unsubscribe", Symbol: symbol})
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func SubscribeToAllActiveSymbols() {
	var activeSymbols []entities.Symbol
	database.DBConn.Where("active = true").Find(&activeSymbols)

	for _, symbol := range activeSymbols {
		subscribe(symbol.Symbol)
	}
}
