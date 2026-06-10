// Package config manages application configuration settings, including environment
// variables, asset directories, database paths, and cryptographic keys.
package config

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/antonlearn/go-final-project/internal/settings"
	"github.com/antonlearn/go-final-project/pkg/hash"
)

// DatabaseConfig holds database-specific attributes including file names,
// resolved absolute paths, and an active connection pool instance.
type DatabaseConfig struct {
	FileName string
	Path     string
	Connect  *sql.DB
}

// EnvVars captures environment-specific configurations needed to run the
// HTTP server, locate static web assets, and authenticate incoming requests.
type EnvVars struct {
	DB               DatabaseConfig
	WebDirPath       string
	Port             string
	ExpectedPassword string
}

// Config encapsulates the entire core application state and runtime constants.
type Config struct {
	Envs         EnvVars
	ExpectedHash string
	JwtKey       []byte
	MaxNumTasks  int
}

// App is the globally accessible configuration instance populated during application startup.
var App Config

// LoadConfig reads environmental flags, resolves absolute runtime paths for static
// assets and databases, computes password hashes, and populates the global App context.
func LoadConfig() {
	password := settings.GetPassword()

	// Compute the SHA-256 hash of the configured password for secure verification.
	var passwordHash string
	if password != "" {
		passwordHash = hash.HashPassword(password)
	}

	// Resolve the application root directory relative to the current executable.
	rootPath := "."
	if exePath, err := os.Executable(); err == nil {
		rootPath = filepath.Dir(exePath)
	}

	// Fallback to the current working directory if web directory assets cannot be found relative to the executable.
	if _, err := os.Stat(filepath.Join(rootPath, settings.GetWebDir())); os.IsNotExist(err) {
		if wd, errWD := os.Getwd(); errWD == nil {
			rootPath = wd
		}
	}

	// Resolve the absolute database path.
	dbFile := settings.GetDBFileName()
	var dbPath string
	if filepath.IsAbs(dbFile) {
		dbPath = dbFile
	} else {
		dbPath = filepath.Join(rootPath, dbFile)
	}

	// Populate the global shared configuration state.
	App = Config{
		Envs: EnvVars{
			DB: DatabaseConfig{
				FileName: dbFile,
				Path:     dbPath,
			},
			WebDirPath:       filepath.Join(rootPath, settings.GetWebDir()),
			Port:             settings.GetPort(),
			ExpectedPassword: password,
		},
		ExpectedHash: passwordHash,
		JwtKey:       []byte(settings.GetSecret()),
		MaxNumTasks:  settings.MaxNumTasks,
	}
}
