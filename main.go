package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"ticker-backend/controllers"
	"ticker-backend/database"
	"ticker-backend/migrations"
)

func main() {
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	fmt.Println("Starting up db connection...")
	database.InitDb()
	migrations.Migrate()

	fmt.Printf("Starting up api on port %s\n", port)

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

	log.Fatal(http.ListenAndServe(":"+port, r))
}
