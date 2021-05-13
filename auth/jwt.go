package auth

import (
	"fmt"
	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secrety-secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func GetToken(claims map[string]interface{}) (string, bool) {
	_, tokenString, err := TokenAuth.Encode(claims)
	if err != nil {
		return "", false
	}

	return tokenString, true
}
