// Package handlers provides HTTP handlers and middleware for the application.
package handlers

import (
	"net/http"

	"github.com/antonlearn/go-final-project/pkg/config"
)

// reloadHomePageHandler returns a handler that serves the frontend files
// from the configured web directory. It also resets the auth cookie
// to ensure the user sees the login page if needed.
func reloadHomePageHandler() http.Handler {
	return resetCookieMiddlewareHandler(
		http.FileServer(http.Dir(config.Config.WebDirPath)),
	)
}
