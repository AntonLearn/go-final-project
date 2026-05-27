package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"

	"github.com/antonlearn/go-final-project/pkg"
)

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
	if _, err := pkg.DB.Exec(initCommand); err != nil {
		return err
	}
	pkg.Logger.Println("Scheduler table and its index file have been successfully created")
	return nil
}

func OpenDB() error {
	// Checking exists of db
	var dbNotExists bool
	if _, err := os.Stat(pkg.DBFile); err != nil {
		dbNotExists = true
	}
	// Opening db
	var err error
	if pkg.DB, err = sql.Open("sqlite", pkg.DBFile); err != nil {
		return err
	}
	if dbNotExists {
		// Creating table and index
		if err := dbInit(); err != nil {
			return err
		}
	}
	pkg.Logger.Printf("Database %s was opened successfully\n", pkg.DBFile)
	return nil
}
