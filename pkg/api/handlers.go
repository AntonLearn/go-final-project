package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antonlearn/go-final-project/pkg"
	"github.com/antonlearn/go-final-project/pkg/db"
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
		now, err = time.Parse(pkg.DateFormatTemplateYYYYMMDD, nowStr)
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
	pkg.Logger.Printf("Next task date %s has been successfully generated and sent by server\n", nextDate)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.GetTasks(search)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, map[string]any{"tasks": tasks})
	pkg.Logger.Printf("List of upcoming tasks %v has been created and sent by server\n", tasks)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeErrorJSON(w, http.StatusBadRequest, "request parameters: no task ID specified")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, emptyMap)
		pkg.Logger.Println("Task was removed from list and processed by server as completed")
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = db.UpdateDateTask(nextDate, id)
		if err != nil {
			writeErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, emptyMap)
		pkg.Logger.Printf("Task was processed by server as completed and its date was changed to new %s\n", nextDate)
	}
}
