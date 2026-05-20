package main

import (
	"log"
	"os"

	"github.com/antonlearn/go-final-project/pkg/db"
	"github.com/antonlearn/go-final-project/pkg/server"
)

func main() {
	// Creating new logger
	logger := log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Opening db
	if err := db.OpenDB(); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Database is open and ready for use")
	defer db.Db.Close()
	// Creating new server with logger
	server := server.NewServer(logger)
	// Starting this server
	if err := server.HTTPServer.ListenAndServe(); err != nil {
		logger.Fatal("Error starting server:", err)
	}
	logger.Printf("Server starting on http://localhost:%s\n", server.HTTPServer.Addr)
}
