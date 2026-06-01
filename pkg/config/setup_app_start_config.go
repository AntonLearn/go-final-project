// Package config contains application configuration and initialization logic.
package config

import (
	"github.com/antonlearn/go-final-project/pkg/hash"
	"github.com/antonlearn/go-final-project/pkg/settings"
)

// SetupAppStartConfig initializes derived configuration values
// after environment variables have been loaded.
func SetupAppStartConfig() {
	// Generate JWT signing key from secret
	Config.JwtKey = []byte(settings.Secret)

	// Hash the expected password for secure comparison
	Config.ExpectedHash = hash.HashPassword(Config.ExpectedPassword)

	// Set maximum number of tasks from settings
	Config.MaxNumTasks = settings.MaxNumTasks
}
