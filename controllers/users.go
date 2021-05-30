package controllers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strings"
	"ticker-backend/auth"
	"ticker-backend/entities"
	"ticker-backend/models"
	"time"
)

type UsersResource struct {
	socketInterrupt chan models.SocketInterrupt
}

func (ur UsersResource) Routes(socketInterrupt chan models.SocketInterrupt) chi.Router {
	router := chi.NewRouter()

	ur.socketInterrupt = socketInterrupt

	//public
	router.Post("/login", ur.Login)
	router.Post("/register", ur.Register)

	//need auth
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(auth.TokenAuth))
		router.Use(auth.Authenticator)

		router.Get("/symbols", ur.GetUserSymbols)
		router.Post("/symbols", ur.AddUserSymbol)

		router.Route("/symbols/{id}", func(router chi.Router) {
			router.Use(UserSymbolCtx)
			router.Delete("/", ur.DeleteUserSymbol)
			router.Put("/", ur.UpdateUserSymbol)
		})
	})

	return router
}

func UserSymbolCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ur UsersResource) Register(w http.ResponseWriter, r *http.Request) {
	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	username := strings.ToLower(strings.TrimSpace(r.Header.Get("username")))
	password := r.Header.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Missing username/password", http.StatusBadRequest)
		return
	}

	if result := db.Find(&entities.User{}, "username = ?", username); result.RowsAffected != 0 {
		http.Error(w, "Username taken", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newUser := entities.User{
		Base:        entities.Base{},
		Username:    username,
		Password:    string(passwordHash),
		UserSymbols: nil,
		LastSeen:    time.Now(),
	}

	db.Create(&newUser)

	token, ok := auth.GetToken(map[string]interface{}{"user_id": newUser.ID.String()})
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(token)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ur UsersResource) Login(w http.ResponseWriter, r *http.Request) {
	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	username := strings.ToLower(strings.TrimSpace(r.Header.Get("username")))
	password := r.Header.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Missing username/password", http.StatusBadRequest)
		return
	}

	var user entities.User

	if result := db.Find(&user, "username = ?", username); result.RowsAffected == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token, ok := auth.GetToken(map[string]interface{}{"user_id": user.ID.String()})
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(token)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ur UsersResource) GetUserSymbols(w http.ResponseWriter, r *http.Request) {
	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user entities.User
	result := db.First(&user, "id = ?", userId)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	//Get the UserSymbols and preload the attached Symbol
	err = db.Model(&user).Preload("Symbol").Association("UserSymbols").Find(&user.UserSymbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//set user symbols to active and subscribe to them if not active
	for _, userSymbol := range user.UserSymbols {
		if userSymbol.Symbol.Active == false {
			//set symbol active
			db.Model(&entities.Symbol{}).Where("id = ?", userSymbol.SymbolID.String()).Update("active", true)

			ur.socketInterrupt <- models.SocketInterrupt{
				InterruptType: "subscribe",
				Symbol:        userSymbol.Symbol.Symbol,
			}
		}
	}

	jsonUserSymbols, ok := ToJson(user.UserSymbols, w)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonUserSymbols); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (ur UsersResource) AddUserSymbol(w http.ResponseWriter, r *http.Request) {
	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	var symbols []models.NewUserSymbol

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(requestBody, &symbols); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuidUserId, err := uuid.FromString(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, symbol := range symbols {
		uuidSymbolId, err := uuid.FromString(symbol.SymbolId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := db.Find(&entities.UserSymbol{}, "user_id = ? and symbol_id = ?", uuidUserId.String(), uuidSymbolId.String())
		if result.RowsAffected != 0 {
			continue //skip since it already exists
		}

		newUserSymbol := entities.UserSymbol{
			UserID:   uuidUserId,
			SymbolID: uuidSymbolId,
			Amount:   symbol.Amount,
		}

		db.Create(&newUserSymbol)

		//set symbol active
		db.Model(&entities.Symbol{}).Where("id = ?", uuidSymbolId.String()).Update("active", true)

		dbSymbol := entities.Symbol{Base: entities.Base{ID: uuidSymbolId}}
		db.First(&dbSymbol)

		//subscribe to symbol
		ur.socketInterrupt <- models.SocketInterrupt{
			InterruptType: "subscribe",
			Symbol:        dbSymbol.Symbol,
		}
	}
}

func (ur UsersResource) DeleteUserSymbol(w http.ResponseWriter, r *http.Request) {
	symbolId := r.Context().Value("id").(string)

	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.Delete(&entities.UserSymbol{}, "user_id = ? and symbol_id = ?", userId, symbolId)
}

func (ur UsersResource) UpdateUserSymbol(w http.ResponseWriter, r *http.Request) {
	symbolId := r.Context().Value("id").(string)

	db, ok := GetDb(w, r)
	if !ok {
		return
	}

	var userSymbol models.NewUserSymbol

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(requestBody, &userSymbol); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := auth.GetUserIdFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.Model(&entities.UserSymbol{}).Where("user_id = ? and symbol_id = ?", userId, symbolId).Update("amount", userSymbol.Amount)
}
