// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/pkg/format"
	"github.com/antonlearn/go-final-project/pkg/logger"
	"github.com/antonlearn/go-final-project/pkg/nextdate"
)

// nextDayHandler calculates and returns the next occurrence date
// based on the provided start date and repetition rule.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var (
		now time.Time
		err error
	)

	// Use current date if 'now' parameter is not provided
	if nowStr == "" {
		now = time.Now().Truncate(24 * time.Hour)
	} else {
		now, err = time.Parse(format.DateFormatTemplateYYYYMMDD, nowStr)
		if err != nil {
			logger.Errorf("Invalid 'now' date format: %s", nowStr)
			writeErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	nextDate, err := nextdate.NextDate(now, dstart, repeat)
	if err != nil {
		logger.Errorf("Failed to calculate next date. now=%s, date=%s, repeat=%s: %v", nowStr, dstart, repeat, err)
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))

	logger.Infof("Next task date %s has been successfully generated and sent by server", nextDate)
}
