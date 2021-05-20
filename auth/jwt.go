package auth

import (
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"ticker-backend/database"
	"ticker-backend/entities"
	"time"
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

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		//check if user still exists
		db := database.DBConn
		var user entities.User

		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		userId := claims["user_id"].(string)

		result := db.Find(&user, "id = ?", userId)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), 500)
			return
		}

		//user doesn't exist
		if result.RowsAffected == 0 {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		//update user LastSeen
		user.LastSeen = time.Now()
		db.Save(&user)

		// Token is authenticated and user exists, pass it through
		next.ServeHTTP(w, r)
	})
}
