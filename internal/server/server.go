// Package server implements HTTP server configuration, timeout settings,
// and operational lifecycle management for the network engine.
package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/antonlearn/go-final-project/internal/config"
	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/internal/handlers"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// httpErrorWriter is an adapter struct that implements the io.Writer interface
// to route standard library error strings directly to the custom logger.
type httpErrorWriter struct {
	logger *logger.Logger
}

// Write intercepts raw byte sequences emitted by the standard server, strips trailing
// whitespace or newlines, and forwards the cleaned output to the application's structured log engine.
func (w httpErrorWriter) Write(p []byte) (n int, err error) {
	cleanMsg := strings.TrimSpace(string(p))
	w.logger.Errorf("HTTP server internal error: %s", cleanMsg)
	return len(p), nil
}

// Server wraps the underlying HTTP server implementation alongside its dedicated logging runtime context.
type Server struct {
	Log        *logger.Logger
	HTTPServer http.Server
}

// NewServer builds, calibrates, and outputs an initialized HTTP server execution context,
// maps routing definitions, redirects operational server logs, and enforces standard timeout constraints.
func NewServer(cfg *config.Config, store *db.Store, appLogger *logger.Logger) *Server {
	mux := http.NewServeMux()

	// Register all application routes and inject dependencies into the central handler router.
	handlers.InitHandlers(mux, store, cfg, appLogger)

	// Intercept standard error streams via the internal writer adapter layer.
	stdErrorLog := log.New(httpErrorWriter{logger: appLogger}, "", 0)

	return &Server{
		Log: appLogger,
		HTTPServer: http.Server{
			Addr:         ":" + cfg.Envs.Port,
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  30 * time.Second,
			ErrorLog:     stdErrorLog,
		},
	}
}
