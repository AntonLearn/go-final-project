package main

import (
	"log"
	"os"

	"github.com/antonlearn/go-final-project/pkg"
	"github.com/antonlearn/go-final-project/pkg/db"
	"github.com/antonlearn/go-final-project/pkg/server"
)

func main() {
	// Creating new logger
	pkg.Logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
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
