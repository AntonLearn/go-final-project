package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"

	"github.com/antonlearn/go-final-project/tests"
)

var Db *sql.DB

const initCommand = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
	title CHAR(255) NOT NULL DEFAULT "",
    comment TEXT DEFAULT "",
    repeat CHAR(128) DEFAULT ""
);
CREATE INDEX idx_date ON scheduler(date);
`

func dbInit() error {
	if _, err := Db.Exec(initCommand); err != nil {
		return err
	}
	return nil
}

func OpenDB() error {
	// Set db file
	dbFile := tests.DBFile[3:]
	if dbFileEnv := os.Getenv("TODO_DBFILE"); dbFileEnv != "" {
		dbFile = dbFileEnv
	}
	// Checking exists of db
	var dbNotExists bool
	if _, err := os.Stat(dbFile); err != nil {
		dbNotExists = true
	}
	// Opening db
	var err error
	if Db, err = sql.Open("sqlite", dbFile); err != nil {
		return err
	}
	if dbNotExists {
		// Creating table and index
		if err := dbInit(); err != nil {
			return err
		}
	}
	return nil
}
