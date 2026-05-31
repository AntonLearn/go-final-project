// Package handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/antonlearn/go-final-project/pkg/config"
)

var emptyMap = make(map[string]any, 0)

func writeErrorJSON(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": errorMessage}); err != nil {
		errorsString := fmt.Sprintf("%s: %s", errorMessage, err.Error())
		http.Error(w, errorsString, http.StatusInternalServerError)
		config.Config.Logger.Println(errorsString)
		return
	}
	config.Config.Logger.Println(errorMessage)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encoded, err := json.Marshal(data)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = w.Write(encoded)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}
