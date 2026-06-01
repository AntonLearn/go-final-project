// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/logger"
	"github.com/antonlearn/go-final-project/pkg/nextdate"
)

// taskDoneHandler marks a task as completed.
// If the task has no repeat rule, it is deleted.
// If it has a repeat rule, its date is moved to the next occurrence.
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeErrorJSON(w, http.StatusBadRequest, "request parameters: no task ID specified")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Errorf("Invalid task ID format: %s", idStr)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		logger.Errorf("Failed to get task with ID %d: %v", id, err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Repeat == "" {
		// One-time task — delete after completion
		err = db.DeleteTask(id)
		if err != nil {
			logger.Errorf("Failed to delete completed one-time task %d: %v", id, err)
			writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, emptyMap)
		logger.Infof("Task %d was removed from list and processed as completed", id)
	} else {
		// Recurring task — calculate next date
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			logger.Errorf("Failed to calculate next date for task %d: %v", id, err)
			writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = db.UpdateDateTask(nextDate, id)
		if err != nil {
			logger.Errorf("Failed to update date for task %d: %v", id, err)
			writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, emptyMap)
		logger.Infof("Task %d was processed as completed and its date was changed to %s", id, nextDate)
	}
}
