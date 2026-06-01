// Package main
package main

import (
	"context"
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
	logWriter, err := logger.SetupLogger()
	if err != nil {
		config.Config.Logger.Println(err)
		return
	}
	defer logWriter.Close()
	if err = config.ReadEnvApp(); err != nil {
		config.Config.Logger.Println(err)
		return
	}
	config.Config.Logger.Println("Application started successfully")
	// Opening db
	if err = db.OpenDB(); err != nil {
		config.Config.Logger.Println(err)
		return
	}
	config.Config.Logger.Printf("Database %s is ready for use", config.Config.DBFileName)
	defer config.Config.DBConnect.Close()
	config.SetupAppStartConfig()
	// Creating new server with logger
	server := server.NewServer()
	mux, ok := server.HTTPServer.Handler.(*http.ServeMux)
	if !ok {
		config.Config.Logger.Println("Expected *http.ServeMux, got something else")
		return
	}
	handlers.InitHandlers(mux)
	// Starting this server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	readyChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)
	baseAddress := "http://localhost" + server.HTTPServer.Addr
	go func() {
		config.Config.Logger.Printf("Server is trying to start on %s\n", baseAddress)
		ln, err := net.Listen("tcp", server.HTTPServer.Addr)
		if err != nil {
			config.Config.Logger.Printf("Failed to listen on %s: %s\n", server.HTTPServer.Addr, err.Error())
			errChan <- err
			return
		}
		readyChan <- struct{}{}
		config.Config.Logger.Println("Server has started listening successfully!")
		if err := server.HTTPServer.Serve(ln); err != nil && err != http.ErrServerClosed {
			config.Config.Logger.Printf("http.Serve returned error: %s\n", err.Error())
			errChan <- err
		} else {
			config.Config.Logger.Println("http.Serve finished with non error or graceful shutdown)")
			errChan <- nil
		}
	}()
	select {
	case <-readyChan:
		config.Config.Logger.Printf("Server is running on %s\n", baseAddress)
	case serverErr := <-errChan:
		if serverErr != nil {
			config.Config.Logger.Printf("Failed to start server: %s\n", serverErr.Error())
			return
		}
	case <-stop:
		config.Config.Logger.Println("Shutdown signal received")
		config.Config.Logger.Println("Initiating graceful shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.HTTPServer.Shutdown(ctx); err != nil {
			config.Config.Logger.Printf("Failed to shutdown server gracefully: %s\n", err.Error())
			config.Config.Logger.Println("Forcing application termination due to failed shutdown")
			os.Exit(1)
		} else {
			config.Config.Logger.Println("Server stopped gracefully")
		}
		config.Config.Logger.Println("Application terminated successfully")
	}
}
