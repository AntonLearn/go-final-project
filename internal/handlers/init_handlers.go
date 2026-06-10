// Package handlers provides HTTP routing, middleware implementation, and endpoint
// lifecycle management for the task scheduler service.
package handlers

import (
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/internal/config"
	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// Handler encapsulates core operational dependencies, including persistence stores,
// environmental configurations, and shared logging tools required across all endpoints.
type Handler struct {
	store  *db.Store
	cfg    *config.Config
	logger *logger.Logger
}

// NewHandler constructs and returns a fully initialized Handler instance via explicit dependency injection.
func NewHandler(store *db.Store, cfg *config.Config, appLogger *logger.Logger) *Handler {
	return &Handler{
		store:  store,
		cfg:    cfg,
		logger: appLogger,
	}
}

// InitHandlers instantiates the centralized router context, maps internal paths to their
// corresponding handler functions, and builds structured middleware chains.
func InitHandlers(mux *http.ServeMux, store *db.Store, cfg *config.Config, appLogger *logger.Logger) {
	h := NewHandler(store, cfg, appLogger)

	// Public routes
	mux.HandleFunc("/health", h.healthCheckHandler)

	// Authentication routes
	mux.HandleFunc("POST /api/signin", h.signinHandler)
	mux.HandleFunc("GET /api/signout", h.signoutHandler)

	// Web frontend routes
	mux.Handle("/", h.reloadHomePageHandler())

	// API routes
	mux.HandleFunc("GET /api/nextdate", h.nextDayHandler)

	// Protected routes (require authentication)
	mux.HandleFunc("GET /api/tasks", h.authMiddlewareHandler(h.tasksHandler))
	mux.HandleFunc("POST /api/task/done", h.authMiddlewareHandler(h.taskDoneHandler))

	// Task CRUD operations (protected)
	mux.HandleFunc("POST /api/task", h.authMiddlewareHandler(h.addTaskHandler))
	mux.HandleFunc("GET /api/task", h.authMiddlewareHandler(h.getTaskHandler))
	mux.HandleFunc("PUT /api/task", h.authMiddlewareHandler(h.updateTaskHandler))
	mux.HandleFunc("DELETE /api/task", h.authMiddlewareHandler(h.deleteTaskHandler))

	h.logger.Info("Handler initialization completed successfully")
}

// healthCheckHandler responds with a standard payload indicating the operational status of the service.
func (h *Handler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := `{"status": "healthy", "timestamp": "` + time.Now().UTC().Format(time.RFC3339) + `"}`
	w.Write([]byte(response))

	h.logger.Infof("Health check called from %s", r.RemoteAddr)
}
