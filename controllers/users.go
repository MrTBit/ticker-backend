package controllers

//import (
//	"context"
//	"github.com/go-chi/chi"
//	"github.com/go-chi/jwtauth"
//	"net/http"
//	"ticker-backend/auth"
//)
//
//type UsersResource struct{}
//
//func (ur UsersResource) Routes() chi.Router {
//	router := chi.NewRouter()
//
//	//public
//	router.Post("/login", ur.Login)
//	router.Post("/register", ur.Register)
//
//	//need auth
//	router.Group(func(router chi.Router){
//		router.Use(jwtauth.Verifier(auth.TokenAuth))
//		router.Use(jwtauth.Authenticator)
//
//		router.Get("/symbols", ur.GetUserSymbols)
//		router.Post("/symbols", ur.AddUserSymbol)
//
//		router.Route("/symbols/{id}", func(router chi.Router) {
//			router.Use(UserSymbolCtx)
//			router.Delete("/", ur.DeleteUserSymbol)
//			router.Put("/", ur.ModifyUserSymbol)
//		})
//	})
//
//	return router
//}
//
//func UserSymbolCtx(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
