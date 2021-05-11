package database

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"time"
)

var DBConn *gorm.DB

func InitDb() {
	var err error

	dsn := "host=<host goes here> port=5432 user=<user goes here>> dbname=<db name here> password=<password goes here>"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "data.", // schema name
		}})

	if err != nil {
		panic("failed to connect to db" + err.Error())
	}
}

func SetDBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeoutContext, _ := context.WithTimeout(context.Background(), time.Second*5)
		ctx := context.WithValue(r.Context(), "DB", DBConn.WithContext(timeoutContext))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
