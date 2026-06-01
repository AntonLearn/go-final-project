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

func main() {
	// Initialize logger
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

	// Open database connection
	if err = db.OpenDB(); err != nil {
		logger.Errorf("Failed to open database: %v", err)
		return
	}
	logger.Infof("Database %s is ready for use", config.Config.DBFileName)
	defer config.Config.DBConnect.Close()

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

	// Graceful Shutdown Setup

	baseAddress := "http://localhost" + srv.HTTPServer.Addr

	// Create cancel function (context itself is not needed here)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Start server in a separate goroutine
	go func() {
		logger.Infof("Server is trying to start on %s", baseAddress)

		ln, err := net.Listen("tcp", srv.HTTPServer.Addr)
		if err != nil {
			logger.Errorf("Failed to listen on %s: %v", srv.HTTPServer.Addr, err)
			cancel()
			return
		}

		logger.Infof("Server started successfully on %s", baseAddress)

		if err := srv.HTTPServer.Serve(ln); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-stop
	logger.Info("Shutdown signal received. Initiating graceful shutdown...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.HTTPServer.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Graceful shutdown failed: %v", err)
		// Force close if graceful shutdown fails
		if closeErr := srv.HTTPServer.Close(); closeErr != nil {
			logger.Errorf("Force close failed: %v", closeErr)
		}
	} else {
		logger.Info("Server stopped gracefully")
	}

	logger.Info("Application terminated successfully")
}
