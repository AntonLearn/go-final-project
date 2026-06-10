// Package handlers implements the HTTP request routing, user authentication,
// and session termination endpoints for the task scheduler service.
package handlers

import (
	"net/http"
)

// signoutHandler invalidates the client session by clearing the active authentication
// token cookie and redirecting the client back to the root application path.
func (h *Handler) signoutHandler(w http.ResponseWriter, r *http.Request) {
	h.resetCookieToken(w)
	h.logger.Info("Redirection to login page completed successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
