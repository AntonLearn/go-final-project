// Package handlers
package handlers

import (
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func resetCookieMiddlewareHandler(handler http.Handler) http.Handler {
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
