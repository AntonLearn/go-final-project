// Package db provides database operations and initialization
// for the task scheduler using SQLite.
package db

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
