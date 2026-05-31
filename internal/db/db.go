// Package db
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"

	"github.com/antonlearn/go-final-project/pkg/config"
)

func dbInit() error {
	if _, err := config.Config.DBConnect.Exec(initCommand); err != nil {
		return err
	}
	config.Config.Logger.Println("Scheduler table and its index file have been successfully created")
	return nil
}

func OpenDB() error {
	// Checking exists of db
	var dbNotExists bool
	if _, err := os.Stat(config.Config.DBPath); err != nil {
		dbNotExists = true
	}
	// Opening db
	var err error
	config.Config.DBConnect, err = sql.Open("sqlite", config.Config.DBPath)
	if err != nil {
		return err
	}
	if dbNotExists {
		// Creating table and index
		if err := dbInit(); err != nil {
			return err
		}
	}
	config.Config.Logger.Printf("Database %s was opened successfully\n", config.Config.DBFileName)
	return nil
}
