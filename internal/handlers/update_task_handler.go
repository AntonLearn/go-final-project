// Package handlers implements the HTTP request routing, request processing logic,
// and resource CRUD operations for the task scheduler service.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/antonlearn/go-final-project/internal/model"
	"github.com/antonlearn/go-final-project/internal/nextdate"
)

// updateTaskHandler reads a task update payload from the request body, executes business logic validations,
// and applies the modified parameters to the existing persisted task record.
func (h *Handler) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var task model.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		h.logger.Errorf("Failed to unmarshal task update request: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Title == "" {
		h.logger.Warn("Update attempt with empty task title")
		h.writeErrorJSON(w, http.StatusBadRequest, "Invalid task title: it's required and cannot be empty")
		return
	}

	// Validate modifications against deadline rules and recurrence schedule formatting.
	if err := nextdate.CheckDate(&task, h.logger); err != nil {
		h.logger.Errorf("Date validation failed for task update: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Apply data adjustments to the corresponding record within the persistent storage layer.
	err = h.store.UpdateTask(&task)
	if err != nil {
		h.logger.Errorf("Failed to update task %s in database: %v", task.ID, err)
		h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, emptyMap)
	h.logger.Infof("Task %s has been successfully updated", task.ID)
}
