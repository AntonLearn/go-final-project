// Package server
package server

import (
	"log"
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/pkg"
	"github.com/antonlearn/go-final-project/pkg/api"
)

type Server struct {
	Log        *log.Logger
	HTTPServer http.Server
}

func NewServer(l *log.Logger) *Server {
	// Creating router
	mux := http.NewServeMux()
	api.InitHandlers(mux)
	// Returning server instance
	return &Server{
		Log: l,
		HTTPServer: http.Server{
			Addr:         ":" + pkg.Port,
			Handler:      mux,
			ErrorLog:     l,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}
