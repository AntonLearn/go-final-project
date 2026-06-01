// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/logger"
	"github.com/antonlearn/go-final-project/pkg/nextdate"
)

// updateTaskHandler updates an existing task with new data from the request body.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Failed to read request body: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		logger.Errorf("Failed to unmarshal task update request: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Title == "" {
		logger.Warn("Update attempt with empty task title")
		writeErrorJSON(w, http.StatusBadRequest, "Invalid task title: it's required and cannot be empty")
		return
	}

	// Validate and adjust date according to repeat rule
	if err := nextdate.CheckDate(&task); err != nil {
		logger.Errorf("Date validation failed for task update: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		logger.Errorf("Failed to update task %s in database: %v", task.ID, err)
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, emptyMap)
	logger.Infof("Task %s has been successfully updated", task.ID)
}
