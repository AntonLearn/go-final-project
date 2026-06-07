// Package config manages application configuration.
package config

import (
	"database/sql"
)

// AppConfig holds all global application settings.

type DatabaseConfig struct {
	FileName string
	Path     string
	Connect  *sql.DB
}

type EnvVars struct {
	DB               DatabaseConfig
	WebDirPath       string
	Port             string
	ExpectedPassword string
}

type AppConfig struct {
	Envs         EnvVars
	ExpectedHash string
	JwtKey       []byte
	MaxNumTasks  int
}

// Config is the global application configuration instance.
var (
	DB     DatabaseConfig
	Envs   EnvVars
	Config AppConfig
)
