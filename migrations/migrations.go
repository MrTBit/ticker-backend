package migrations

import (
	"ticker-backend/database"
	"ticker-backend/models"
)

func Migrate() {
	database.DBConn.AutoMigrate(models.Symbol{}, models.User{})
}
