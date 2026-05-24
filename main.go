package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antonlearn/go-final-project/pkg"
	"github.com/antonlearn/go-final-project/pkg/db"
	"github.com/antonlearn/go-final-project/pkg/server"
)

func main() {
	logFile, err := SetupLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	pkg.Logger.Println("Application started successfully")
	// Opening db
	if err := db.OpenDB(); err != nil {
		pkg.Logger.Fatal(err)
	}
	pkg.Logger.Println("Database is ready for use")
	defer db.Db.Close()
	// Creating new server with logger
	server := server.NewServer(pkg.Logger)
	// Starting this server
	pkg.Logger.Printf("Server starting on http://localhost:%s\n", server.HTTPServer.Addr)
	if err := server.HTTPServer.ListenAndServe(); err != nil {
		pkg.Logger.Fatal("Error starting server:", err)
	}
}

func SetupLogger() (*os.File, error) {
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
