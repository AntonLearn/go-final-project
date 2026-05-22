package api

import (
	"encoding/json"
	"net/http"

	"github.com/antonlearn/go-final-project/pkg"
)

func writeErrorJSON(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
	pkg.Logger.Println(errorMessage)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encoded, err := json.Marshal(data)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Write(encoded)
}
