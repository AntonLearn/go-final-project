// Package handlers provides helper functions for writing JSON responses
// and handling errors in HTTP handlers.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/antonlearn/go-final-project/pkg/logger"
)

var emptyMap = make(map[string]any, 0)

// writeErrorJSON writes an error response in JSON format
func writeErrorJSON(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": errorMessage}); err != nil {
		errorsString := fmt.Sprintf("%s: %s", errorMessage, err.Error())

		// Log the failure to encode error response
		logger.Errorf("Failed to encode error response: %v. Original error: %s", err, errorMessage)

		http.Error(w, errorsString, http.StatusInternalServerError)
		return
	}

	// Log the error that was sent to the client
	logger.Errorf("Error response sent: %s (status: %d)", errorMessage, status)
}

// writeJSON writes any data as a JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	encoded, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Failed to marshal JSON data: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(encoded)
	if err != nil {
		logger.Errorf("Failed to write JSON response to client: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}
