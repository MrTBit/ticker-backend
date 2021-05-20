package tickers

import (
	"ticker-backend/database"
	"ticker-backend/entities"
	"ticker-backend/socket"
	"time"
)

type ActiveUserTicker struct{}

func (aut ActiveUserTicker) Start() {
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			<-ticker.C
			//look for inactive users (haven't called in in the last minute) and unsubscribe from any symbols only they used
			type symbolToUnsubscribe struct {
				Id     string
				Symbol string
			}
			var symbolsToUnsubscribe []symbolToUnsubscribe
			database.DBConn.Model(&entities.Symbol{}).Select("data.symbols.id, data.symbols.symbol").
				Joins("join data.user_symbols us on symbols.id = us.symbol_id").
				Joins("join data.users u on us.user_id = u.id").
				Where("last_seen < CURRENT_TIMESTAMP - INTERVAL '1 minute' AND us.symbol_id not in (SELECT symbol_id FROM data.user_symbols JOIN data.users u ON data.user_symbols.user_id = u.id WHERE u.last_seen >= CURRENT_TIMESTAMP - INTERVAL '1 minute')").
				Group("data.symbols.id").
				Find(&symbolsToUnsubscribe)

			for _, symbol := range symbolsToUnsubscribe {
				socket.Unsubscribe(symbol.Symbol)
				//TODO: convert to update where id in [] to hit db less
				database.DBConn.Model(&entities.Symbol{}).Where("id = ?", symbol.Id).Update("active", false)
			}
		}
	}()
}
