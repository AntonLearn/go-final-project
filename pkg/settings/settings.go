// Package settings contains default configuration constants for the application.
package settings

const (
	// Database settings
	DBFileName = "scheduler.db"

	// Web frontend directory
	WebDir = "web"

	// Server settings
	Port = "7540"

	// Authentication settings
	ExpectedPassword = "" // empty by default (can be overridden via TODO_PASSWORD env)

	// JWT secret key (should be changed in production!)
	Secret = "your-secret-key-here"

	// Application limits
	MaxNumTasks = 50
)
