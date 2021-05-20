package tickers

import "ticker-backend/models"

func Init(socketInterrupt chan models.SocketInterrupt) {
	ActiveUserTicker{}.Start(socketInterrupt)
}
