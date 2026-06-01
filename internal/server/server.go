// Package server provides HTTP server configuration and initialization.
package server

import (
	"log"
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// Server wraps the HTTP server and related components.
type Server struct {
	Log        *log.Logger
	HTTPServer http.Server
}

// NewServer creates and configures a new HTTP server instance
// with appropriate timeouts and middleware-ready handler.
func NewServer() *Server {
	// Create the main request router
	mux := http.NewServeMux()

	return &Server{
		// Using the new leveled logger (Info level) for general logging
		Log: logger.InfoLogger,
		HTTPServer: http.Server{
			Addr:    ":" + config.Config.Port,
			Handler: mux,
			// Use ErrorLogger for server's internal error logging
			ErrorLog:     logger.ErrorLogger,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}
