// Package handlers implements the HTTP request routing, request processing logic,
// and collection retrieval endpoints for the task scheduler service.
package handlers

import (
	"net/http"
)

// tasksHandler extracts an optional search expression from the query string parameters
// and serves a JSON array containing matching task entities retrieved from the store.
func (h *Handler) tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	// Fetch the filtered or un-filtered task collection via the storage layer.
	tasks, err := h.store.GetTasks(search)
	if err != nil {
		h.logger.Errorf("Failed to get tasks with search query '%s': %v", search, err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, map[string]any{"tasks": tasks})

	if search != "" {
		h.logger.Infof("List of tasks (search: '%s') sent. Total tasks: %d", search, len(tasks))
	} else {
		h.logger.Infof("List of all upcoming tasks sent. Total tasks: %d", len(tasks))
	}
}
