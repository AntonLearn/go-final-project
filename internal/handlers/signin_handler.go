// Package handlers
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if config.Config.ExpectedPassword == "" {
		config.Config.Logger.Println("Empty password. Dummy-token was successfully created and sent by server")
		config.Config.Logger.Println("Login completed successfully")
		writeJSON(w, map[string]string{"token": "dummy-token"})
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	var request struct {
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if request.Password != config.Config.ExpectedPassword {
		writeErrorJSON(w, http.StatusUnauthorized, "Invalid password")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"password_hash": config.Config.ExpectedHash})
	tokenString, err := token.SignedString(config.Config.JwtKey)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, map[string]string{"token": tokenString})
	config.Config.Logger.Println("Correct password. Token was successfully created and sent by server")
	config.Config.Logger.Println("Login completed successfully")
}
