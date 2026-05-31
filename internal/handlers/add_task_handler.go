// Package handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/nextdate"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if task.Title == "" {
		writeErrorJSON(w, http.StatusBadRequest, "Invalid task title: it's required and cannot be empty")
		return
	}
	if err := nextdate.CheckDate(&task); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	var id int64
	id, err = db.AddTask(&task)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
	config.Config.Logger.Println("Task was successfully added and sent by server")
}
