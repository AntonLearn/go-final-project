// Package settings provides hardcoded configuration fallbacks and dynamic environment
// variable initialization helpers for the task scheduler infrastructure.
package settings

import "os"

const (
	// Default configuration parameters applied when explicit runtime flags are absent.
	defaultDBFileName = "scheduler.db"
	defaultWebDir     = "web"
	defaultPort       = "7540"
	defaultSecret     = "your-secret-key-here"
	defaultPasword    = "password"

	// MaxNumTasks caps the maximum capacity boundary for individual task query iterations.
	MaxNumTasks = 50

	// LogFileStdout directs target streams for core logger routines (File, Stdout, Both).
	LogFileStdout = "File"
)

// GetPort extracts the network listening interface boundary from the TODO_PORT
// environment context, defaulting to standard port 7540 if unassigned.
func GetPort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return defaultPort
}

// GetPassword determines password via the TODO_PASSWORD environment context,
// returning the default password if empty.
func GetPassword() string {
	if password := os.Getenv("TODO_PASSWORD"); password != "" {
		return password
	}
	return defaultPasword
}

// GetDBFileName determines the structural storage file destination path via the
// TODO_DBFILE environment context, returning the default sqlite database path if empty.
func GetDBFileName() string {
	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		return dbFile
	}
	return defaultDBFileName
}

// GetSecret fetches the cryptographic token signing string from the TODO_SECRET
// environment context, reverting to a hardcoded baseline key phrase if absent.
func GetSecret() string {
	if secret := os.Getenv("TODO_SECRET"); secret != "" {
		return secret
	}
	return defaultSecret
}

// GetWebDir extracts the workspace directory tracking compiled web assets
// from the TODO_WEB_DIR environment context, falling back to the default distribution root.
func GetWebDir() string {
	if webDir := os.Getenv("TODO_WEB_DIR"); webDir != "" {
		return webDir
	}
	return defaultWebDir
}
