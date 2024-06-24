package authentication

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"

	"awesomeProject/pkg/utils"
)

func IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetTokenFromCookie(r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusUnauthorized)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    claims["username"].(string),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "userID",
			Value:    claims["userID"].(string),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		})

		next.ServeHTTP(w, r)
	}
}

func GetTokenFromCookie(r *http.Request) (jwt.MapClaims, error) {
	token, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	claims, err := utils.VerifyJWT(token.Value)
	if err != nil {
		return nil, err
	}

	expirationDate, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if expirationDate.Before(time.Now()) {
		return nil, err
	}

	return claims, nil
}
