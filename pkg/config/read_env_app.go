// Package config
package config

import (
	"os"
	"path/filepath"

	"github.com/antonlearn/go-final-project/pkg/settings"
)

func ReadEnvApp() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	rootPath := filepath.Dir(exePath)
	Config.DBFileName = settings.DBFileName
	if dbFileNameEnv := os.Getenv("TODO_DBFILE"); dbFileNameEnv != "" {
		Config.DBFileName = dbFileNameEnv
	}
	Config.DBPath = filepath.Join(rootPath, Config.DBFileName)
	Config.WebDirPath = filepath.Join(rootPath, settings.WebDir)
	Config.Port = settings.Port
	if portEnv := os.Getenv("TODO_PORT"); portEnv != "" {
		Config.Port = portEnv
	}
	Config.ExpectedPassword = settings.ExpectedPassword
	if expectedPasswordEnv := os.Getenv("TODO_PASSWORD"); expectedPasswordEnv != "" {
		Config.ExpectedPassword = expectedPasswordEnv
	}
	return nil
}
