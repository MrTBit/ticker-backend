package controllers

import (
	"context"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ticker-backend/auth"
	"ticker-backend/models"
	"time"
)

type UsersResource struct{}

func (ur UsersResource) Routes() chi.Router {
	router := chi.NewRouter()

	//public
	router.Post("/login", ur.Login)
	router.Post("/register", ur.Register)

	////need auth
	//router.Group(func(router chi.Router){
	//	router.Use(jwtauth.Verifier(auth.TokenAuth))
	//	router.Use(jwtauth.Authenticator)
	//
	//	router.Get("/symbols", ur.GetUserSymbols)
	//	router.Post("/symbols", ur.AddUserSymbol)
	//
	//	router.Route("/symbols/{id}", func(router chi.Router) {
	//		router.Use(UserSymbolCtx)
	//		router.Delete("/", ur.DeleteUserSymbol)
	//		router.Put("/", ur.ModifyUserSymbol)
	//	})
	//})

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

	username := r.Header.Get("username")
	password := r.Header.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Missing username/password", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Base:       models.Base{},
		Username:   username,
		Password:   string(passwordHash),
		UserSymbol: nil,
		LastSeen:   time.Now(),
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

	username := r.Header.Get("username")
	password := r.Header.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Missing username/password", http.StatusBadRequest)
		return
	}

	var user models.User
	db.First(&user, "username = ?", username)

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
