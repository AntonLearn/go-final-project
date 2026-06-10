// Package handlers implements the HTTP request routing, request processing logic,
// and resource CRUD operations for the task scheduler service.
package handlers

import (
	"net/http"
	"strconv"
)

// deleteTaskHandler extracts the task identifier from the query parameters and triggers its permanent removal from the persistence layer.
func (h *Handler) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// Delete the task record via the storage layer.
	err = h.store.DeleteTask(id)
	if err != nil {
		h.logger.Errorf("Failed to delete task %d: %v", id, err)
		h.writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, emptyMap)
	h.logger.Infof("Task %d was successfully deleted by server", id)
}
