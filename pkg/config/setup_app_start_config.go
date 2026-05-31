// Package config
package config

import (
	"github.com/antonlearn/go-final-project/pkg/hash"
	"github.com/antonlearn/go-final-project/pkg/settings"
)

func SetupAppStartConfig() {
	Config.JwtKey = []byte(settings.Secret)
	Config.ExpectedHash = hash.HashPassword(Config.ExpectedPassword)
	Config.MaxNumTasks = settings.MaxNumTasks
}
