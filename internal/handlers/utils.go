// Package handlers implements the HTTP request routing, response serialization,
// and utility helper functions for the task scheduler service.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// emptyMap represents a pre-allocated, reusable empty JSON object
// to prevent unnecessary dynamic memory allocations across empty responses.
var emptyMap = make(map[string]any, 0)

// writeErrorJSON serializes and transmits a structured error message
// payload along with the corresponding HTTP status code.
func (h *Handler) writeErrorJSON(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": errorMessage}); err != nil {
		errorsString := fmt.Sprintf("%s: %s", errorMessage, err.Error())
		h.logger.Errorf("Failed to encode error response: %v", err)
		http.Error(w, errorsString, http.StatusInternalServerError)
		return
	}

	h.logger.Errorf("Error response sent: %s (status: %d)", errorMessage, status)
}

// writeJSON serializes any given interface data structures into a JSON payload
// and writes the resulting byte array directly to the HTTP response stream.
func (h *Handler) writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	encoded, err := json.Marshal(data)
	if err != nil {
		h.logger.Errorf("Failed to marshal JSON data: %v", err)
		h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(encoded)
	if err != nil {
		h.logger.Errorf("Failed to write JSON response: %v", err)
		h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}
