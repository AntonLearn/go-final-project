// Package db manages database operations, schema initializations, and structural constants
// for the task scheduler using SQLite as the storage backend.
package db

// initCommand contains the SQL schema definition required to initialize the scheduler table and its associated indexes.
const initCommand = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title CHAR(255) NOT NULL DEFAULT "",
    comment TEXT DEFAULT "",
    repeat CHAR(128) DEFAULT ""
);

-- Index for fast date-based queries
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`
