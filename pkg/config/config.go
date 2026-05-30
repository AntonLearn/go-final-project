// Package config
package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/antonlearn/go-final-project/pkg/hash"
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

const (
	dbFileName       = "scheduler.db"
	webDir           = "web"
	port             = "7540"
	expectedPassword = ""
	secret           = "your-secret-key-here"
	maxNumTasks      = 50
)

func SetConfigureApp() {
	Config.JwtKey = []byte(secret)
	Config.ExpectedHash = hash.HashPassword(Config.ExpectedPassword)
	Config.MaxNumTasks = maxNumTasks
}

func ReadEnvApp() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	rootPath := filepath.Dir(exePath)
	Config.DBFileName = dbFileName
	if dbFileNameEnv := os.Getenv("TODO_DBFILE"); dbFileNameEnv != "" {
		Config.DBFileName = dbFileNameEnv
	}
	Config.DBPath = filepath.Join(rootPath, Config.DBFileName)
	Config.WebDirPath = filepath.Join(rootPath, webDir)
	Config.Port = port
	if portEnv := os.Getenv("TODO_PORT"); portEnv != "" {
		Config.Port = portEnv
	}
	Config.ExpectedPassword = expectedPassword
	if expectedPasswordEnv := os.Getenv("TODO_PASSWORD"); expectedPasswordEnv != "" {
		Config.ExpectedPassword = expectedPasswordEnv
	}
	return nil
}
