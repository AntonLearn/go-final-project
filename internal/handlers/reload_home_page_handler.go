// Package handlers implements the HTTP request routing, static file distribution,
// and session management middleware for the task scheduler service.
package handlers

import (
	"net/http"
)

// reloadHomePageHandler configures and returns an HTTP file server instance wrapped
// in cookie remediation middleware to distribute static frontend assets.
func (h *Handler) reloadHomePageHandler() http.Handler {
	return h.resetCookieMiddlewareHandler(
		http.FileServer(http.Dir(h.cfg.Envs.WebDirPath)),
	)
}
