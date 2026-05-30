// Package api
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/format"
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
	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
	config.Config.Logger.Printf("Next task date %s has been successfully generated and sent by server\n", nextDate)
}

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
		config.Config.Logger.Println("Task was removed from list and processed by server as completed")
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
		config.Config.Logger.Printf("Task was processed by server as completed and its date was changed to new %s\n", nextDate)
	}
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if config.Config.ExpectedPassword == "" {
		config.Config.Logger.Println("Empty password. Dummy-token was successfully created and sent by server")
		config.Config.Logger.Println("Login completed successfully")
		writeJSON(w, map[string]string{"token": "dummy-token"})
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	var request struct {
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if request.Password != config.Config.ExpectedPassword {
		writeErrorJSON(w, http.StatusUnauthorized, "Invalid password")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"password_hash": config.Config.ExpectedHash})
	tokenString, err := token.SignedString(config.Config.JwtKey)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, map[string]string{"token": tokenString})
	config.Config.Logger.Println("Correct password. Token was successfully created and sent by server")
	config.Config.Logger.Println("Login completed successfully")
}

func signoutHandler(w http.ResponseWriter, r *http.Request) {
	resetCookieToken(w)
	config.Config.Logger.Println("Redirection to login page completed successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func authMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.Config.ExpectedPassword == "" {
			config.Config.Logger.Println("Authentication completed successfully")
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
			func(token *jwt.Token) (any, error) {
				return config.Config.JwtKey, nil
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
		if storedHash != config.Config.ExpectedHash {
			writeErrorJSON(w, http.StatusUnauthorized, "authentication required: password verification failed. Stored hash does not match current password")
			return
		}
		config.Config.Logger.Println("Authentication completed successfully")
		handler(w, r)
	}
}

func reloadHomePageHandler() http.Handler {
	return resetCookieMiddleware(http.FileServer(http.Dir(config.Config.WebDirPath)))
}

func resetCookieMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resetCookieToken(w)
		config.Config.Logger.Println("Login page has been reloaded successfully")
		handler.ServeHTTP(w, r)
	})
}

func resetCookieToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		SameSite: http.SameSiteDefaultMode,
	})
	config.Config.Logger.Println("Token in cookies was deleted successfully")
	config.Config.Logger.Println("Logout completed successfully")
}
