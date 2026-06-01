// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/logger"
	"github.com/antonlearn/go-final-project/pkg/nextdate"
)

// addTaskHandler creates a new task from the request body.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Failed to read request body: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		logger.Errorf("Failed to unmarshal task creation request: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Title == "" {
		logger.Warn("Attempt to create task with empty title")
		writeErrorJSON(w, http.StatusBadRequest, "Invalid task title: it's required and cannot be empty")
		return
	}

	// Validate and adjust date according to repeat rule
	if err := nextdate.CheckDate(&task); err != nil {
		logger.Errorf("Date validation failed for new task: %v", err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		logger.Errorf("Failed to add task to database: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
	logger.Infof("Task was successfully added with ID %d", id)
}
