// Package handlers implements the HTTP request routing, request processing logic,
// and resource CRUD operations for the task scheduler service.
package handlers

import (
	"net/http"
	"strconv"
)

// getTaskHandler retrieves a single task using the unique identifier provided in the query parameters.
func (h *Handler) getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch the task record from the storage layer.
	task, err := h.store.GetTask(id)
	if err != nil {
		h.logger.Errorf("Failed to get task %d: %v", id, err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, task)
	h.logger.Infof("Task %d was successfully retrieved and sent by server", id)
}
