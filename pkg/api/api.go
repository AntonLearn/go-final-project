package api

import (
	"net/http"
)

var WebDir = "./web" // Directory with frontend files

func InitHandlers(mux *http.ServeMux) {
	// Creating handler for root path and all subpaths of webDir
	mux.Handle("/", http.FileServer(http.Dir(WebDir)))
	// Creating handler for api/nextdate
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	// Creating handler for api/task
	mux.HandleFunc("/api/task", taskHandler)
	// Creating handler for api/tasks
	mux.HandleFunc("/api/tasks", tasksHandler)
}
