// Package db provides database operations for the task scheduler
// using SQLite as the storage backend.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// PackageLogger defines the localized logging interface required by the db package,
// decoupling it from external logging implementations.
type PackageLogger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Info(args ...any)
}

// Store encapsulates the database connection pool, logger instance, and operational constraints.
type Store struct {
	db          *sql.DB
	logger      PackageLogger
	maxNumTasks int
}

// NewStore opens a SQLite database connection, ensures the destination directory exists,
// and initializes the schema if the database file is newly created.
func NewStore(dbPath string, maxNumTasks int, appLogger PackageLogger) (*Store, error) {
	dbNotExists := false
	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			dbNotExists = true
		} else {
			appLogger.Errorf("Failed to check database file status: %v", err)
			return nil, fmt.Errorf("failed to check database file status: %w", err)
		}
	}

	// Ensure the parent directory for the database file exists.
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		appLogger.Errorf("Failed to create directories for database path %s: %v", dbDir, err)
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		appLogger.Errorf("Failed to open SQLite database at %s: %v", dbPath, err)
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	store := &Store{
		db:          db,
		logger:      appLogger,
		maxNumTasks: maxNumTasks,
	}

	if dbNotExists {
		if err := store.dbInit(); err != nil {
			db.Close()
			return nil, err
		}
	}

	appLogger.Infof("Database %s was opened successfully", filepath.Base(dbPath))
	return store, nil
}

// dbInit executes the initialization command to set up the necessary tables and indexes.
func (s *Store) dbInit() error {
	if _, err := s.db.Exec(initCommand); err != nil {
		s.logger.Errorf("Failed to initialize database schema: %v", err)
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	s.logger.Info("Scheduler table and indexes have been successfully created")
	return nil
}

// Close gracefully terminates the underlying database connection pool.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
