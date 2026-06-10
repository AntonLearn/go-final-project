// Package handlers implements the HTTP request routing, request processing logic,
// and utility calculation endpoints for the task scheduler service.
package handlers

import (
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/internal/nextdate"
	"github.com/antonlearn/go-final-project/pkg/format"
)

// nextDayHandler handles calculations for recurrent task execution timelines,
// parsing baseline dates, target scopes, and recurrence rule strings to return a plain-text date response.
func (h *Handler) nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var (
		now time.Time
		err error
	)

	if nowStr == "" {
		// Default to the current system date truncated to midnight if no simulation baseline is specified.
		now = time.Now().Truncate(24 * time.Hour)
	} else {
		now, err = time.Parse(format.YYYYMMDD, nowStr)
		if err != nil {
			h.logger.Errorf("Invalid 'now' date format: %s", nowStr)
			h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Delegate scheduling computation to the underlying rule engine.
	nextDate, err := nextdate.NextDate(now, dstart, repeat, h.logger)
	if err != nil {
		h.logger.Errorf("Failed to calculate next date: %v", err)
		h.writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
	h.logger.Infof("Next task date %s has been successfully generated", nextDate)
}
