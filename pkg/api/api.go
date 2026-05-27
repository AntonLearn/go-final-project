// Package api
package api

import (
	"net/http"

	"github.com/antonlearn/go-final-project/pkg"
)

func InitHandlers(mux *http.ServeMux) {
	const webDir = "./web" // Directory with frontend files
	// Creating handler for /api/signin
	mux.HandleFunc("POST /api/signin", signinHandler)
	// Creating handler for /api/signout
	mux.HandleFunc("GET /api/signout", signoutHandler)
	// Creating handler for root path and all subpaths of webDir
	mux.Handle("GET /", reloadHomePageHandler(webDir))
	// Creating handler for api/nextdate
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	// Creating handler for api/tasks
	mux.HandleFunc("GET /api/tasks", authMiddleware(tasksHandler))
	// Creating handler for api/task/done
	mux.HandleFunc("POST /api/task/done", authMiddleware(taskDoneHandler))
	// Creating handlers for api/task
	mux.HandleFunc("POST /api/task", authMiddleware(addTaskHandler))
	mux.HandleFunc("GET /api/task", authMiddleware(getTaskHandler))
	mux.HandleFunc("PUT /api/task", authMiddleware(updateTaskHandler))
	mux.HandleFunc("DELETE /api/task", authMiddleware(deleteTaskHandler))
	pkg.Logger.Println("Handler initialization completed successfully")
}
