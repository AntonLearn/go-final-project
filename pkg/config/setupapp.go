// Package config manages application configuration, including
// environment variable parsing and path setup.
package config

import (
	"os"
	"path/filepath"

	"github.com/antonlearn/go-final-project/pkg/hash"
	"github.com/antonlearn/go-final-project/pkg/settings"
)

// SetupApp initializes the application's configuration values
func SetupApp() error {
	// Read environment variables into a structure for storing them.
	// If any environment variables are not defined, set their default values.
	if err := Envs.readEnvApp(); err != nil {
		return err
	}
	Config.Envs = Envs
	Config.setConfig()
	return nil
}

// ReadEnvApp reads configuration from environment variables and sets default values.
func (env *EnvVars) readEnvApp() error {
	// Get the directory where the executable is located
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	rootPath := filepath.Dir(exePath)

	// Read basic database settings from environment variables
	DB.setDBBaseSettings(rootPath)
	Envs.DB = DB

	// Web frontend directory
	env.WebDirPath = filepath.Join(rootPath, settings.WebDir)

	// Server port
	env.Port = settings.Port
	if portEnv := os.Getenv("TODO_PORT"); portEnv != "" {
		env.Port = portEnv
	}

	// Password for authentication
	env.ExpectedPassword = settings.ExpectedPassword
	if expectedPasswordEnv := os.Getenv("TODO_PASSWORD"); expectedPasswordEnv != "" {
		env.ExpectedPassword = expectedPasswordEnv
	}

	return nil
}

// setDBBaseSettings reads basic database settings from environment variables
func (db *DatabaseConfig) setDBBaseSettings(rootPath string) {
	// Database base settings
	db.FileName = settings.DBFileName
	if dbFileNameEnv := os.Getenv("TODO_DBFILE"); dbFileNameEnv != "" {
		db.FileName = dbFileNameEnv
	}
	db.Path = filepath.Join(rootPath, db.FileName)
}

// setConfig sets the values of the application's configuration structure.
func (cfg *AppConfig) setConfig() {
	// Generate JWT signing key from secret
	cfg.JwtKey = []byte(settings.Secret)

	// Hash the expected password for secure comparison
	cfg.ExpectedHash = hash.HashPassword(cfg.Envs.ExpectedPassword)

	// Set maximum number of tasks from settings
	cfg.MaxNumTasks = settings.MaxNumTasks
}
