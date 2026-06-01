// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// authMiddlewareHandler checks JWT token validity and password hash match.
// It protects routes that require authentication.
func authMiddlewareHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If no password is set in config, skip authentication (development mode)
		if config.Config.ExpectedPassword == "" {
			logger.Info("Authentication skipped - no password configured")
			handler(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			logger.Warnf("No auth token cookie found from %s", r.RemoteAddr)
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return config.Config.JwtKey, nil
		})

		if err != nil || !token.Valid {
			logger.Warnf("Invalid or expired JWT token from %s", r.RemoteAddr)
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}

		value, exists := claims["password_hash"]
		if !exists {
			logger.Warn("JWT token missing password_hash claim")
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password_hash claim is missing")
			return
		}

		storedHash, ok := value.(string)
		if !ok {
			logger.Warnf("password_hash claim has wrong type: %T", value)
			writeErrorJSON(w, http.StatusUnauthorized,
				fmt.Sprintf("authentication required: password_hash has unexpected type: %T but password_hash must be a string", value))
			return
		}

		if storedHash != config.Config.ExpectedHash {
			logger.Warn("Password hash verification failed - token does not match current password")
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password verification failed. Stored hash does not match current password")
			return
		}

		logger.Info("Authentication completed successfully")
		handler(w, r)
	}
}
