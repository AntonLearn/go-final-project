// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"

	"github.com/antonlearn/go-final-project/pkg/logger"
)

// signoutHandler logs out the user by resetting the auth cookie
// and redirects them to the login page.
func signoutHandler(w http.ResponseWriter, r *http.Request) {
	resetCookieToken(w)

	logger.Info("Redirection to login page completed successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
