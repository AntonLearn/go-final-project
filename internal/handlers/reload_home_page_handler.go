// Package handlers
package handlers

import (
	"net/http"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func reloadHomePageHandler() http.Handler {
	return resetCookieMiddlewareHandler(http.FileServer(http.Dir(config.Config.WebDirPath)))
}
