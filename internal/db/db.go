// Package db provides database operations for the task scheduler
// using SQLite as the storage backend.
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite" // Pure Go SQLite driver

	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

// dbInit creates the necessary tables and indexes if they don't exist.
func dbInit() error {
	if _, err := config.Config.DBConnect.Exec(initCommand); err != nil {
		logger.Errorf("Failed to initialize database schema: %v", err)
		return err
	}

	logger.Info("Scheduler table and indexes have been successfully created")
	return nil
}

// OpenDB opens the SQLite database and initializes the schema if needed.
func OpenDB() error {
	// Check if database file exists
	dbNotExists := false
	if _, err := os.Stat(config.Config.DBPath); err != nil {
		if os.IsNotExist(err) {
			dbNotExists = true
		} else {
			logger.Errorf("Failed to check database file status: %v", err)
			return err
		}
	}

	// Open database connection
	var err error
	config.Config.DBConnect, err = sql.Open("sqlite", config.Config.DBPath)
	if err != nil {
		logger.Errorf("Failed to open SQLite database at %s: %v", config.Config.DBPath, err)
		return err
	}

	// Create tables if this is a new database
	if dbNotExists {
		if err := dbInit(); err != nil {
			return err
		}
	}

	logger.Infof("Database %s was opened successfully", config.Config.DBFileName)
	return nil
}
