// Package main
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antonlearn/go-final-project/pkg"
	"github.com/antonlearn/go-final-project/pkg/api"
	"github.com/antonlearn/go-final-project/pkg/db"
	"github.com/antonlearn/go-final-project/pkg/server"
	"github.com/antonlearn/go-final-project/tests"
)

func main() {
	logFile, err := setupLogger()
	if err != nil {
		pkg.Logger.Println(err)
		return
	}
	defer logFile.Close()
	readConfig()
	pkg.Logger.Println("Application started successfully")
	// Opening db
	if err := db.OpenDB(); err != nil {
		pkg.Logger.Fatal(err)
	}
	pkg.Logger.Println("Database is ready for use")
	defer pkg.DB.Close()
	// Creating new server with logger
	server := server.NewServer(pkg.Logger)
	// Starting this server
	pkg.Logger.Printf("Server starting on http://localhost:%s\n", server.HTTPServer.Addr)
	if err := server.HTTPServer.ListenAndServe(); err != nil {
		pkg.Logger.Fatal("Error starting server:", err)
	}
}

func setupLogger() (*os.File, error) {
	file, err := os.Create(generationLocalFileName(".log"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", file.Name(), err)
	}
	// Creating new logger
	pkg.Logger = log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return file, nil
}

func generationLocalFileName(ext string) string {
	timeStamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	fileName := "app_" + timeStamp + ext
	return fileName
}

func readConfig() {
	pkg.Port = fmt.Sprintf("%d", tests.Port)
	if portEnv := os.Getenv("TODO_PORT"); portEnv != "" {
		pkg.Port = portEnv
	}
	pkg.DBFile = tests.DBFile[3:]
	if dbFileEnv := os.Getenv("TODO_DBFILE"); dbFileEnv != "" {
		pkg.DBFile = dbFileEnv
	}
	pkg.ExpectedPassword = ""
	if expectedPasswordEnv := os.Getenv("TODO_PASSWORD"); expectedPasswordEnv != "" {
		pkg.ExpectedPassword = expectedPasswordEnv
	}
	api.HashPassword()
}
