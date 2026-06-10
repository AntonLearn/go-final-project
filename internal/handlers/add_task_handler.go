// Package handlers implements the HTTP request routing and payload validation logic
// for the task scheduler service, ensuring secure and decoupled data flow.
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/antonlearn/go-final-project/internal/model"
	"github.com/antonlearn/go-final-project/internal/nextdate"
)

// addTaskHandler processes incoming JSON payloads to validate, parse, and persist a new task entity.
func (h *Handler) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var task model.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		h.logger.Errorf("Failed to unmarshal task creation request: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Title == "" {
		h.logger.Warn("Attempt to create task with empty title")
		h.writeErrorJSON(w, http.StatusBadRequest, "Invalid task title: it's required and cannot be empty")
		return
	}

	// Validate execution dates and process recurrence schedules if provided.
	if err := nextdate.CheckDate(&task, h.logger); err != nil {
		h.logger.Errorf("Date validation failed for new task: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Persist the verified task record via the storage layer.
	id, err := h.store.AddTask(&task)
	if err != nil {
		h.logger.Errorf("Failed to add task to database: %v", err)
		h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
	h.logger.Infof("Task was successfully added with ID %d", id)
}
