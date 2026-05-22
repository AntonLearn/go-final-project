package api

import (
	"net/http"

	"github.com/antonlearn/go-final-project/pkg"
)

var WebDir = "./web" // Directory with frontend files

func InitHandlers(mux *http.ServeMux) {
	// Creating handler for root path and all subpaths of webDir
	mux.Handle("/", http.FileServer(http.Dir(WebDir)))
	// Creating handler for api/nextdate
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	// Creating handler for api/tasks
	mux.HandleFunc("/api/tasks", tasksHandler)
	// Creating handler for api/task/done
	mux.HandleFunc("/api/task/done", taskDoneHandler)
	// Creating handlers for api/task
	mux.HandleFunc("POST /api/task", addTaskHandler)
	mux.HandleFunc("GET /api/task", getTaskHandler)
	mux.HandleFunc("PUT /api/task", updateTaskHandler)
	mux.HandleFunc("DELETE /api/task", deleteTaskHandler)
	pkg.Logger.Println("Handler initialization completed successfully")
}
