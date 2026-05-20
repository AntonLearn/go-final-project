package tests

import (
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

const (
	initCommand = `
		CREATE TABLE scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date CHAR(8) NOT NULL DEFAULT "",
		title CHAR(255) NOT NULL DEFAULT "",
    	comment TEXT DEFAULT "",
    	repeat CHAR(128) DEFAULT ""
	);
	CREATE INDEX idx_date ON scheduler(date);
	`
	countCommand  = `SELECT count(id) FROM scheduler`
	insertCommand = `INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (?, 'Todo', 'Комментарий', '')`
	selectCommand = `SELECT * FROM scheduler WHERE id=?`
	deleteCommand = `DELETE FROM scheduler WHERE id = ?`
)

func count(db *sqlx.DB) (int, error) {
	var count int
	return count, db.Get(&count, countCommand)
}

func dbInit(db *sqlx.DB) error {
	if _, err := db.Exec(initCommand); err != nil {
		return err
	}
	return nil
}

func openDB(t *testing.T) *sqlx.DB {
	dbfile := DBFile
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbfile = envFile
	}
	// Checking exists of db
	var dbNotExists bool
	if _, err := os.Stat(dbfile); err != nil {
		dbNotExists = true
	}
	db, err := sqlx.Connect("sqlite", dbfile)
	assert.NoError(t, err)
	if dbNotExists {
		err := dbInit(db)
		assert.NoError(t, err)
	}
	return db
}

func TestDB(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	before, err := count(db)
	assert.NoError(t, err)

	today := time.Now().Format(`20060102`)

	res, err := db.Exec(insertCommand, today)
	assert.NoError(t, err)

	id, err := res.LastInsertId()

	var task Task
	err = db.Get(&task, selectCommand, id)
	assert.NoError(t, err)
	assert.Equal(t, id, task.ID)
	assert.Equal(t, `Todo`, task.Title)
	assert.Equal(t, `Комментарий`, task.Comment)

	_, err = db.Exec(deleteCommand, id)
	assert.NoError(t, err)

	after, err := count(db)
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}
