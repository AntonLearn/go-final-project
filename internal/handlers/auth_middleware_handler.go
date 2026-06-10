// Package handlers implements the HTTP request routing, request processing logic,
// and security middleware for the task scheduler service.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// authMiddlewareHandler wraps an http.HandlerFunc to enforce JSON Web Token (JWT)
// authentication and verify claims matching the configured password hash.
func (h *Handler) authMiddlewareHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Bypass authentication if no system-wide security password has been configured.
		if h.cfg.Envs.ExpectedPassword == "" {
			h.logger.Info("Authentication skipped - no password configured")
			handler(w, r)
			return
		}

		// Extract the JWT session token from the incoming client cookies.
		cookie, err := r.Cookie("token")
		if err != nil {
			h.logger.Warnf("No auth token cookie found from %s", r.RemoteAddr)
			h.writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}

		// Parse and validate the token signature against the configured signing key.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return h.cfg.JwtKey, nil
		})

		if err != nil || !token.Valid {
			h.logger.Warnf("Invalid or expired JWT token from %s", r.RemoteAddr)
			h.writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}

		// Verify the existence of the required password hash claim within the token context.
		value, exists := claims["password_hash"]
		if !exists {
			h.logger.Warn("JWT token missing password_hash claim")
			h.writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password_hash claim is missing")
			return
		}

		storedHash, ok := value.(string)
		if !ok {
			h.logger.Warnf("password_hash claim has wrong type: %T", value)
			h.writeErrorJSON(w, http.StatusUnauthorized, fmt.Sprintf("authentication required: unexpected type: %T", value))
			return
		}

		// Assert that the token's password hash matches the latest active application state hash.
		if storedHash != h.cfg.ExpectedHash {
			h.logger.Warn("Password hash verification failed")
			h.writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password verification failed")
			return
		}

		h.logger.Info("Authentication completed successfully")
		handler(w, r)
	}
}
