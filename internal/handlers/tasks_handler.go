// Package handlers
package handlers

import (
	"net/http"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/config"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.GetTasks(search)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, map[string]any{"tasks": tasks})
	config.Config.Logger.Printf("List of upcoming tasks %v has been created and sent by server\n", tasks)
}
