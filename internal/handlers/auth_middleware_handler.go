// Package handlers
package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func authMiddlewareHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.Config.ExpectedPassword == "" {
			config.Config.Logger.Println("Authentication completed successfully")
			handler(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}
		tokenString := cookie.Value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (any, error) {
				return config.Config.JwtKey, nil
			})
		if err != nil || !token.Valid {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}
		value, exists := claims["password_hash"]
		if !exists {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password_hash claim is missing")
			return
		}
		storedHash, ok := value.(string)
		if !ok {
			writeErrorJSON(w, http.StatusUnauthorized, fmt.Sprintf("authentication required: password_hash has unexpected type: %T but password_hash must be a string", value))
			return
		}
		if storedHash != config.Config.ExpectedHash {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password verification failed. Stored hash does not match current password")
			return
		}
		config.Config.Logger.Println("Authentication completed successfully")
		handler(w, r)
	}
}
