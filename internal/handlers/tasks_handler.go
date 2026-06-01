// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// tasksHandler returns a list of tasks, optionally filtered by search query.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.GetTasks(search)
	if err != nil {
		logger.Errorf("Failed to get tasks with search query '%s': %v", search, err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, map[string]any{"tasks": tasks})

	if search != "" {
		logger.Infof("List of tasks (search: '%s') has been created and sent by server. Total tasks: %d", search, len(tasks))
	} else {
		logger.Infof("List of all upcoming tasks has been created and sent by server. Total tasks: %d", len(tasks))
	}
}
