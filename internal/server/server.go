// Package server
package server

import (
	"log"
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/internal/api"
	"github.com/antonlearn/go-final-project/pkg/config"
)

type Server struct {
	Log        *log.Logger
	HTTPServer http.Server
}

func NewServer() *Server {
	// Creating router
	mux := http.NewServeMux()
	api.InitHandlers(mux)
	// Returning server instance
	return &Server{
		Log: config.Config.Logger,
		HTTPServer: http.Server{
			Addr:         ":" + config.Config.Port,
			Handler:      mux,
			ErrorLog:     config.Config.Logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}
