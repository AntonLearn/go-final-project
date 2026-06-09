// Package main serves as the entry point for the task scheduler service,
// handling configuration loading, infrastructure initialization, and graceful shutdown.
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

	"github.com/antonlearn/go-final-project/internal/config"
	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/internal/server"
	"github.com/antonlearn/go-final-project/internal/settings"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

const failExitCode = 1

func main() {
	if err := run(); err != nil {
		os.Exit(failExitCode)
	}
}

func run() error {
	// Initialize the application logger according to configuration settings (stdout, file, or both).
	appLogger, err := logger.New(settings.LogFileStdout)
	if err != nil {
		fmt.Printf("Failed to setup logger: %v\n", err)
		return err
	}
	defer appLogger.Close()
	appLogger.Info("Logger successfully initialized")
	appLogger.Info("Starting application lifecycle execution...")

	// Load global application configurations.
	appLogger.Info("Loading global application configuration...")
	config.LoadConfig()
	appLogger.Info("Global application configuration successfully loaded")

	appLogger.Info("Initializing application infrastructure...")

	// Initialize the isolated database store with configuration variables.
	store, err := db.NewStore(config.App.Envs.DB.Path, config.App.MaxNumTasks, appLogger)
	if err != nil {
		appLogger.Errorf("Failed to open database: %v", err)
		return err
	}
	appLogger.Infof("Database %s is ready for use", config.App.Envs.DB.FileName)

	defer func() {
		if err := store.Close(); err != nil {
			appLogger.Errorf("Error closing database: %v", err)
		}
	}()

	// Initialize the HTTP server instance with explicit dependency injection.
	srv := server.NewServer(&config.App, store, appLogger)

	// Bind and start the network TCP listener.
	listener, err := net.Listen("tcp", srv.HTTPServer.Addr)
	if err != nil {
		appLogger.Errorf("Failed to bind to address %s: %v", srv.HTTPServer.Addr, err)
		return err
	}
	defer listener.Close()

	// Log successful initialization of the network listener interface.
	appLogger.Infof("Network TCP listener successfully started and binding to address: %s", srv.HTTPServer.Addr)

	// Format a clean base address URL for user-friendly initialization logs.
	actualAddr := listener.Addr().String()
	host, port, err := net.SplitHostPort(actualAddr)
	var baseAddress string
	if err == nil {
		if host == "::" || host == "0.0.0.0" || host == "" {
			host = "localhost"
		}
		baseAddress = fmt.Sprintf("http://%s:%s", host, port)
	} else {
		baseAddress = "http://" + actualAddr
	}

	// Set up signal interception to handle termination events gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	serverErrors := make(chan error, 1)

	// Run the HTTP server asynchronously in a separate goroutine.
	go func() {
		if err := srv.HTTPServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Log operational readiness state accurately after the background HTTP server process has been launched.
	appLogger.Infof("Server successfully started and listening on %s", baseAddress)
	appLogger.Info("Application is fully operational and ready to accept traffic")

	// Block execution until an error occurs or a system interrupt signal is received.
	select {
	case err := <-serverErrors:
		appLogger.Errorf("Server encountered a critical runtime error: %v", err)
		return err
	case <-ctx.Done():
		appLogger.Info("Shutdown signal received. Initiating graceful shutdown...")
		stop()
	}

	// Enforce a strict 30-second timeout context for flushing in-flight requests.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.HTTPServer.Shutdown(shutdownCtx); err != nil {
		appLogger.Errorf("Graceful shutdown failed: %v", err)
		if closeErr := srv.HTTPServer.Close(); closeErr != nil {
			appLogger.Errorf("Force close failed: %v", closeErr)
		}
	} else {
		appLogger.Info("Server stopped gracefully")
	}

	appLogger.Info("Application terminated successfully")
	return nil
}
