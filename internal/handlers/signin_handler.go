// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// signinHandler handles user authentication by password and returns a JWT token.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	// If no password is configured, return a dummy token (development mode)
	if config.Config.Envs.ExpectedPassword == "" {
		logger.Info("Empty password. Dummy-token was successfully created and sent by server")
		logger.Info("Login completed successfully")
		writeJSON(w, map[string]string{"token": "dummy-token"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Failed to read request body: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Errorf("Failed to unmarshal signin request: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if request.Password != config.Config.Envs.ExpectedPassword {
		logger.Warnf("Invalid password attempt from %s", r.RemoteAddr)
		writeErrorJSON(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// Create JWT token with password hash claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"password_hash": config.Config.ExpectedHash})

	tokenString, err := token.SignedString(config.Config.JwtKey)
	if err != nil {
		logger.Errorf("Failed to sign JWT token: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, map[string]string{"token": tokenString})

	logger.Info("Correct password. Token was successfully created and sent by server")
	logger.Info("Login completed successfully")
}
