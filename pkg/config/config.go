// Package config
package config

import (
	"database/sql"
	"log"
)

type ConfigApp struct {
	Logger           *log.Logger
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

var Config ConfigApp
