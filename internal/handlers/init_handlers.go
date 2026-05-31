// Package handlers
package handlers

import (
	"net/http"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func InitHandlers(mux *http.ServeMux) {
	// Creating handler for /api/signin
	mux.HandleFunc("POST /api/signin", signinHandler)
	// Creating handler for /api/signout
	mux.HandleFunc("GET /api/signout", signoutHandler)
	// Creating handler for root path and all subpaths of webDir
	mux.Handle("GET /", reloadHomePageHandler())
	// Creating handler for api/nextdate
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	// Creating handler for api/tasks
	mux.HandleFunc("GET /api/tasks", authMiddlewareHandler(tasksHandler))
	// Creating handler for api/task/done
	mux.HandleFunc("POST /api/task/done", authMiddlewareHandler(taskDoneHandler))
	// Creating handlers for api/task
	mux.HandleFunc("POST /api/task", authMiddlewareHandler(addTaskHandler))
	mux.HandleFunc("GET /api/task", authMiddlewareHandler(getTaskHandler))
	mux.HandleFunc("PUT /api/task", authMiddlewareHandler(updateTaskHandler))
	mux.HandleFunc("DELETE /api/task", authMiddlewareHandler(deleteTaskHandler))
	config.Config.Logger.Println("Handler initialization completed successfully")
}
