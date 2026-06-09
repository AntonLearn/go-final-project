// Package handlers implements the HTTP request routing, static file distribution,
// and session management middleware for the task scheduler service.
package handlers

import (
	"net/http"
	"time"
)

// resetCookieMiddlewareHandler intercepts incoming HTTP requests to proactively clear active
// authentication session tokens before passing execution down to the underlying handler.
func (h *Handler) resetCookieMiddlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.resetCookieToken(w)
		h.logger.Info("Login page has been reloaded successfully")
		handler.ServeHTTP(w, r)
	})
}

// resetCookieToken explicitly invalidates the client's session state by issuing an instantly
// expired "token" cookie targeting the root application path.
func (h *Handler) resetCookieToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		SameSite: http.SameSiteDefaultMode,
	})

	h.logger.Info("Token in cookies was deleted successfully")
}
