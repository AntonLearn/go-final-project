// Package handlers implements the HTTP request routing, user authentication,
// and session management endpoints for the task scheduler service.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// signinHandler verifies incoming credentials against the system configuration and generates a signed JWT upon success.
func (h *Handler) signinHandler(w http.ResponseWriter, r *http.Request) {
	// Bypass explicit signature validation if the application is running without a mandatory password restriction.
	if h.cfg.Envs.ExpectedPassword == "" {
		h.logger.Info("Empty password. Dummy-token was successfully created and sent by server")
		h.writeJSON(w, map[string]string{"token": "dummy-token"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.logger.Errorf("Failed to unmarshal signin request: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if request.Password != h.cfg.Envs.ExpectedPassword {
		h.logger.Warnf("Invalid password attempt from %s", r.RemoteAddr)
		h.writeErrorJSON(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// Issue a new token embedding the pre-computed hash value to maintain uniform context maps.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"password_hash": h.cfg.ExpectedHash})
	tokenString, err := token.SignedString(h.cfg.JwtKey)
	if err != nil {
		h.logger.Errorf("Failed to sign JWT token: %v", err)
		h.writeErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeJSON(w, map[string]string{"token": tokenString})
	h.logger.Info("Login completed successfully")
}
