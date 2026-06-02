// Package main is the entry point of the application.
// It initializes all components and starts the HTTP server with graceful shutdown support.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/internal/handlers"
	"github.com/antonlearn/go-final-project/internal/server"
	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

func main() {
	// Initialize the structured logger
	logWriter, err := logger.SetupLogger()
	if err != nil {
		// Fallback to console if logger initialization fails
		fmt.Printf("Failed to setup logger: %v\n", err)
		return
	}
	defer logWriter.Close()

	// Read application configuration from environment variables
	if err = config.ReadEnvApp(); err != nil {
		logger.Errorf("Failed to read environment configuration: %v", err)
		return
	}

	logger.Info("Application started successfully")

	// Open connection to the database
	dbConnect, err := db.OpenDB()
	if err != nil {
		logger.Errorf("Failed to open database: %v", err)
		return
	}
	logger.Infof("Database %s is ready for use", config.Config.DBFileName)

	// Ensure database connection is safety closed on application exit
	defer func(dbConnect *sql.DB) {
		if dbConnect != nil {
			if err := db.Close(); err != nil {
				logger.Errorf("Error closing database: %v", err)
			}
		}
	}(dbConnect)

	config.SetupAppStartConfig()

	// Create HTTP server instance
	srv := server.NewServer()

	// Initialize request handlers
	mux, ok := srv.HTTPServer.Handler.(*http.ServeMux)
	if !ok {
		logger.Error("Expected *http.ServeMux, got something else")
		return
	}
	handlers.InitHandlers(mux)

	baseAddress := "http://localhost" + srv.HTTPServer.Addr

	// Server launch tracking (split stages)

	// Stage 1: Explicitly log the intent/attempt to start the server
	logger.Infof("Attempting to bind and listen on address: %s", baseAddress)

	// Explicitly try to open the network port. If this fails (e.g. port already in use),
	// the application will fail here immediately before starting any background routines
	listener, err := net.Listen("tcp", srv.HTTPServer.Addr)
	if err != nil {
		logger.Errorf("Failed to bind to address %s: %v", srv.HTTPServer.Addr, err)
	}
	defer listener.Close()

	// Stage 2: Explicitly log that the port is successfully bound and the server is actually running
	logger.Infof("Server successfully started and listening on %s", baseAddress)

	// Graceful Shutdown Setup

	// Create a context that listens for termination signals from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Channel to capture critical errors while the server is serving requests
	serverErrors := make(chan error, 1)

	// Start processing active connection requests in a separate goroutine
	go func() {
		// Pass the pre-established listener instead of using ListenAndServe
		if err := srv.HTTPServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Block the main goroutine until an OS signal is caught or a runtime server error occurs
	select {
	case err := <-serverErrors:
		logger.Errorf("Server encountered a critical runtime error: %v", err)
		return
	case <-ctx.Done():
		// OS signal received, proceed to graceful shutdown sequence
		logger.Info("Shutdown signal received. Initiating graceful shutdown...")
		stop() // Stop receiving further signal notifications as early as possible
	}

	// Create a context with a timeout for the graceful shutdown duration
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Attempt to shutdown the server gracefully, allowing active connections to finish
	if err := srv.HTTPServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Graceful shutdown failed: %v", err)
		// Force immediate closure if the graceful shutdown times out
		if closeErr := srv.HTTPServer.Close(); closeErr != nil {
			logger.Errorf("Force close failed: %v", closeErr)
		}
	} else {
		logger.Info("Server stopped gracefully")
	}

	logger.Info("Application terminated successfully")
}
