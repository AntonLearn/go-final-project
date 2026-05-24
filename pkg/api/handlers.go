package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"

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

func signinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	var request struct {
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJSON(w, map[string]string{"token": "dummy-token"})
		return
	}
	if request.Password != expectedPassword {
		writeErrorJSON(w, http.StatusUnauthorized, "Invalid password")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"password_hash": HashPassword(expectedPassword)})
	tokenString, err := token.SignedString(pkg.JwtKey)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, map[string]string{"token": tokenString})
	pkg.Logger.Println("Correct password. Token was successfully created and sent by server")
}

func authMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expectedPassword := os.Getenv("TODO_PASSWORD")
		if expectedPassword == "" {
			handler(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}
		tokenString := cookie.Value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (interface{}, error) {
				return pkg.JwtKey, nil
			})
		if err != nil || !token.Valid {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: "+err.Error())
			return
		}
		value, exists := claims["password_hash"]
		if !exists {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password_hash claim is missing")
			return
		}
		storedHash, ok := value.(string)
		if !ok {
			writeErrorJSON(w, http.StatusUnauthorized, fmt.Sprintf("authentication required: password_hash has unexpected type: %T but password_hash must be a string", value))
			return
		}
		if storedHash != HashPassword(expectedPassword) {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password verification failed. Stored hash does not match current password")
			return
		}
		pkg.Logger.Println("Authentication completed successfully")
		handler(w, r)
	}
}
