// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// getTaskHandler returns a single task by its ID.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeErrorJSON(w, http.StatusBadRequest, "request parameters: no task ID specified")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Errorf("Invalid task ID format: %s", idStr)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		logger.Errorf("Failed to get task %d: %v", id, err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, task)
	logger.Infof("Task %d was successfully retrieved and sent by server", id)
}
