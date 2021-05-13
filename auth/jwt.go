package auth

import (
	"github.com/go-chi/jwtauth"
	"net/http"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secrety-secret"), nil)
}

func GetToken(claims map[string]interface{}) (string, bool) {
	_, tokenString, err := TokenAuth.Encode(claims)
	if err != nil {
		return "", false
	}

	return tokenString, true
}

func GetUserIdFromContext(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", err
	}
	return claims["user_id"].(string), nil
}
