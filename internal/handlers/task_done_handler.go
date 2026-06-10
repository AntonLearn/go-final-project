// Package handlers implements the HTTP request routing, request processing logic,
// and lifecycle management for scheduled tasks.
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antonlearn/go-final-project/internal/nextdate"
)

// taskDoneHandler marks a specific task as completed. If the task is non-recurring,
// it is permanently removed from the store; otherwise, its schedule is advanced to the next calculated runtime.
func (h *Handler) taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.writeErrorJSON(w, http.StatusBadRequest, "request parameters: no task ID specified")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorf("Invalid task ID format: %s", idStr)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Retrieve the task state via the storage layer to evaluate its recurrence settings.
	task, err := h.store.GetTask(id)
	if err != nil {
		h.logger.Errorf("Failed to get task with ID %d: %v", id, err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Repeat == "" {
		// Process one-time tasks by removing them completely from data persistence.
		err = h.store.DeleteTask(id)
		if err != nil {
			h.logger.Errorf("Failed to delete completed one-time task %d: %v", id, err)
			h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		h.writeJSON(w, emptyMap)
		h.logger.Infof("Task %d was removed from list and processed as completed", id)
	} else {
		// Advance recurring tasks to their subsequent calendar milestones.
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat, h.logger)
		if err != nil {
			h.logger.Errorf("Failed to calculate next date for task %d: %v", id, err)
			h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = h.store.UpdateDateTask(nextDate, id)
		if err != nil {
			h.logger.Errorf("Failed to update date for task %d: %v", id, err)
			h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		h.writeJSON(w, emptyMap)
		h.logger.Infof("Task %d processed as completed, next date: %s", id, nextDate)
	}
}
