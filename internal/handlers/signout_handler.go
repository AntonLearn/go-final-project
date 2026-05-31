// Package handlers
package handlers

import (
	"net/http"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func signoutHandler(w http.ResponseWriter, r *http.Request) {
	resetCookieToken(w)
	config.Config.Logger.Println("Redirection to login page completed successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
