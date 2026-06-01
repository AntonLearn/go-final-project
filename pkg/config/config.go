// Package config manages application configuration.
package config

import "database/sql"

// AppConfig holds all global application settings.
type AppConfig struct {
	DBFileName       string
	DBPath           string
	WebDirPath       string
	Port             string
	ExpectedPassword string
	ExpectedHash     string
	JwtKey           []byte
	DBConnect        *sql.DB
	MaxNumTasks      int
}

// Config is the global application configuration instance.
var Config AppConfig
