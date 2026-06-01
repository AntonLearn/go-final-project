// Package config manages application configuration, including
// environment variable parsing and path setup.
package config

import (
	"os"
	"path/filepath"

	"github.com/antonlearn/go-final-project/pkg/settings"
)

// ReadEnvApp reads configuration from environment variables and sets default values.
func ReadEnvApp() error {
	// Get the directory where the executable is located
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	rootPath := filepath.Dir(exePath)

	// Database settings
	Config.DBFileName = settings.DBFileName
	if dbFileNameEnv := os.Getenv("TODO_DBFILE"); dbFileNameEnv != "" {
		Config.DBFileName = dbFileNameEnv
	}
	Config.DBPath = filepath.Join(rootPath, Config.DBFileName)

	// Web frontend directory
	Config.WebDirPath = filepath.Join(rootPath, settings.WebDir)

	// Server port
	Config.Port = settings.Port
	if portEnv := os.Getenv("TODO_PORT"); portEnv != "" {
		Config.Port = portEnv
	}

	// Password for authentication
	Config.ExpectedPassword = settings.ExpectedPassword
	if expectedPasswordEnv := os.Getenv("TODO_PASSWORD"); expectedPasswordEnv != "" {
		Config.ExpectedPassword = expectedPasswordEnv
	}

	return nil
}
