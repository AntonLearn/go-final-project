// Package main
package main

import (
	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/internal/server"
	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

func main() {
	logFileName, err := logger.SetupLogger()
	if err != nil {
		config.Config.Logger.Println(err)
		return
	}
	defer logFileName.Close()
	err = config.ReadEnvApp()
	if err != nil {
		config.Config.Logger.Println(err)
		return
	}
	config.Config.Logger.Println("Application started successfully")
	// Opening db
	err = db.OpenDB()
	if err != nil {
		config.Config.Logger.Println(err)
		return
	}
	config.Config.Logger.Printf("Database %s is ready for use", config.Config.DBFileName)
	defer config.Config.DBConnect.Close()
	config.SetConfigureApp()
	// Creating new server with logger
	server := server.NewServer()
	// Starting this server
	config.Config.Logger.Printf("Server starting on http://localhost:%s\n", server.HTTPServer.Addr)
	if err := server.HTTPServer.ListenAndServe(); err != nil {
		config.Config.Logger.Println("Error starting server:", err)
		return
	}
}
