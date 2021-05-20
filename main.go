package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"ticker-backend/controllers"
	"ticker-backend/database"
	"ticker-backend/models"
	"ticker-backend/socket"
	"ticker-backend/tickers"
)

func main() {
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	log.Println("Starting up db connection...")
	database.InitDb()

	log.Println("Connecting to websocket...")
	socketInterrupt := make(chan models.SocketInterrupt)
	go socket.InitSocket(socketInterrupt)

	log.Println("Starting tickers")
	tickers.Init(socketInterrupt)

	log.Printf("Starting up api on port %s\n", port)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(database.SetDBMiddleware)

	if loggingDisabled := os.Getenv("DISABLE_LOGGING"); loggingDisabled == "" {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello There"))
	})

	r.Mount("/symbols", controllers.SymbolsResource{}.Routes())
	r.Mount("/users", controllers.UsersResource{}.Routes(socketInterrupt))

	log.Fatal(http.ListenAndServe(":"+port, r))
}
