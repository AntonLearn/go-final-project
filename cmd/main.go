// Package main is the entry point of the application.
// It initializes all components and starts the HTTP server with graceful shutdown support.
package main

import (
	"context"
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

// Standard error exit code for operational failures
const failExitCode = 1

func main() {
	// Execute core application logic inside the run abstraction layer to preserve defer semantics.
	if err := run(); err != nil {
		// Standard error exit code for operational failures.
		os.Exit(failExitCode)
	}
}

func run() error {
	// Initialize the structured logger
	logWriter, err := logger.SetupLogger()
	if err != nil {
		// Fallback to console if logger initialization fails
		fmt.Printf("Failed to setup logger: %v\n", err)
		return err
	}
	defer logWriter.Close()

	// Read application configuration from environment variables
	if err = config.SetupApp(); err != nil {
		logger.Errorf("Failed to read environment configuration: %v", err)
		return err
	}

	logger.Info("Initializing application infrastructure...")

	// Open connection to the database
	err = db.OpenDB()
	if err != nil {
		logger.Errorf("Failed to open database: %v", err)
		return err
	}
	logger.Infof("Database %s is ready for use", config.DB.FileName)

	// Automatically executes database cleanup on wrapper termination
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("Error closing database: %v", err)
		}
	}()

	// Create HTTP server instance
	srv := server.NewServer()

	// Initialize request handlers
	mux, ok := srv.HTTPServer.Handler.(*http.ServeMux)
	if !ok {
		logger.Error("Expected *http.ServeMux, got something else")
		return fmt.Errorf("invalid handler type assertion")
	}
	handlers.InitHandlers(mux)

	// Explicitly try to open the network port.
	listener, err := net.Listen("tcp", srv.HTTPServer.Addr)
	if err != nil {
		logger.Errorf("Failed to bind to address %s: %v", srv.HTTPServer.Addr, err)
		return err
	}
	defer listener.Close()

	// Extract the actual address assigned by the Operating System
	actualAddr := listener.Addr().String()

	// Parse host and port to replace generalized interfaces (0.0.0.0 or ::) with localhost for user convenience
	host, port, err := net.SplitHostPort(actualAddr)
	var baseAddress string
	if err == nil {
		if host == "::" || host == "0.0.0.0" || host == "" {
			host = "localhost"
		}
		baseAddress = fmt.Sprintf("http://%s:%s", host, port)
	} else {
		// Fallback to raw address if splitting fails
		baseAddress = "http://" + actualAddr
	}

	// Stage 2: Explicitly log that the port is successfully bound and the server is actually running
	logger.Infof("Server successfully started and listening on %s", baseAddress)
	logger.Info("Application is fully operational and ready to accept traffic")

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
		return err // Non-zero exit cascade via run execution failure
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
	return nil // Clean exit with status code 0
}
