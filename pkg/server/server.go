package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antonlearn/go-final-project/pkg/api"
	"github.com/antonlearn/go-final-project/tests"
)

type Server struct {
	Log        *log.Logger
	HTTPServer http.Server
}

func NewServer(l *log.Logger) *Server {
	// Creating router
	mux := http.NewServeMux()
	api.InitHandlers(mux)
	// Set port for listen
	port := fmt.Sprintf("%d", tests.Port)
	if portEnv := os.Getenv("TODO_PORT"); portEnv != "" {
		port = portEnv
	}
	// Returning server instance
	return &Server{
		Log: l,
		HTTPServer: http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ErrorLog:     l,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}
