package migrations

import (
	"ticker-backend/database"
	"ticker-backend/models"
)

func Migrate() {
	database.DBConn.AutoMigrate(models.Symbol{}, models.User{}, models.UserSymbol{})

	err := database.DBConn.SetupJoinTable(&models.User{}, "UserSymbols", &models.UserSymbol{})

	if err != nil {
		panic("failed to connect to db: " + err.Error())
	}
}
