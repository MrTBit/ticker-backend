package controllers

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"net/http"
	"ticker-backend/auth"
	"ticker-backend/entities"
)

type SymbolsResource struct{}

func (sr SymbolsResource) Routes() chi.Router {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {

		router.Use(jwtauth.Verifier(auth.TokenAuth)) //seek and verify jwt tokens
		router.Use(jwtauth.Authenticator)            //handle valid/invalid tokens -> sends 401 if not valid, otherwise goes through

		router.Get("/", sr.List) //GET /symbols - Read all symbols

		router.Route("/{id}", func(router chi.Router) {
			router.Use(SymbolCtx)      //id context
			router.Get("/", sr.GetOne) //GET /symbols/{id} - read specific symbol
		})
	})

	return router
}

func SymbolCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sr SymbolsResource) List(w http.ResponseWriter, r *http.Request) {
	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	symbolSearch := r.URL.Query().Get("symbol")
	descriptionSearch := r.URL.Query().Get("description")

	var symbols []entities.Symbol

	if symbolSearch != "" || descriptionSearch != "" {
		//fetch based on symbol/description
		db.Where("symbol ilike ?", "%"+symbolSearch+"%").Where("description ilike ?", "%"+descriptionSearch+"%").Find(&symbols)
	} else {
		//fetch all
		db.Find(&symbols)
	}

	symbolsJson, ok := ToJson(symbols, w)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(symbolsJson); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (sr SymbolsResource) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)

	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	var symbol entities.Symbol
	db.Preload("UserSymbols").First(&symbol, "id = ?", id)

	jsonSymbol, ok := ToJson(symbol, w)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonSymbol); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
