// Package handlers
package handlers

import (
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/format"
	"github.com/antonlearn/go-final-project/pkg/nextdate"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")
	var (
		now time.Time
		err error
	)
	if nowStr == "" {
		now = time.Now().Truncate(24 * time.Hour)
	} else {
		now, err = time.Parse(format.DateFormatTemplateYYYYMMDD, nowStr)
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	nextDate, err := nextdate.NextDate(now, dstart, repeat)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
	config.Config.Logger.Printf("Next task date %s has been successfully generated and sent by server\n", nextDate)
}
