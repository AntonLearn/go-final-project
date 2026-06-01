// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/pkg/logger"
)

// InitHandlers registers all application routes and their corresponding handlers.
// It also applies authentication middleware where required.
func InitHandlers(mux *http.ServeMux) {
	// Public routes
	mux.HandleFunc("/health", healthCheckHandler)

	// Authentication routes
	mux.HandleFunc("POST /api/signin", signinHandler)
	mux.HandleFunc("GET /api/signout", signoutHandler)

	// Web frontend routes
	mux.Handle("GET /", reloadHomePageHandler())

	// API routes
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)

	// Protected routes (require authentication)
	mux.HandleFunc("GET /api/tasks", authMiddlewareHandler(tasksHandler))
	mux.HandleFunc("POST /api/task/done", authMiddlewareHandler(taskDoneHandler))

	// Task CRUD operations (protected)
	mux.HandleFunc("POST /api/task", authMiddlewareHandler(addTaskHandler))
	mux.HandleFunc("GET /api/task", authMiddlewareHandler(getTaskHandler))
	mux.HandleFunc("PUT /api/task", authMiddlewareHandler(updateTaskHandler))
	mux.HandleFunc("DELETE /api/task", authMiddlewareHandler(deleteTaskHandler))

	logger.Info("Handler initialization completed successfully")
}

// healthCheckHandler returns the current health status of the application.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := `{"status": "healthy", "timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `"}`
	w.Write([]byte(response))

	logger.Infof("Health check called from %s", r.RemoteAddr)
}
