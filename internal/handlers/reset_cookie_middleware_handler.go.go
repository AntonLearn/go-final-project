// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"
	"time"

	"github.com/antonlearn/go-final-project/pkg/logger"
)

// resetCookieMiddlewareHandler is a middleware that resets the authentication token cookie
// before serving the login page. Useful for forcing re-login.
func resetCookieMiddlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resetCookieToken(w)
		logger.Info("Login page has been reloaded successfully")
		handler.ServeHTTP(w, r)
	})
}

// resetCookieToken removes the authentication token by setting an expired cookie.
// This effectively logs the user out.
func resetCookieToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		SameSite: http.SameSiteDefaultMode,
	})

	logger.Info("Token in cookies was deleted successfully")
	logger.Info("Logout completed successfully")
}
