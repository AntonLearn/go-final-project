package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/antonlearn/go-final-project/pkg"
)

const (
	addTaskCommand = `INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (:date, :title, :comment, :repeat)`
	getTasksLimitBaseCommand      = `SELECT * FROM scheduler ORDER BY date LIMIT :limit`
	getTasksLimitWhereDateCommand = `SELECT * FROM scheduler WHERE date = :date 
		ORDER BY date LIMIT :limit`
	getTasksLimitWhereTitleOrCommentCommand = `SELECT * FROM scheduler 
		WHERE LOWER(title) LIKE LOWER(:search) OR LOWER(comment) LIKE LOWER(:search) ORDER BY date LIMIT :limit`
)

func AddTask(task *Task) (int64, error) {
	result, err := Db.Exec(addTaskCommand, sql.Named("date", task.Date), sql.Named("title",
		task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	var (
		tasks []*Task
		rows  *sql.Rows
		err   error
	)
	search = strings.TrimSpace(search)
	if search != "" {
		searchDate, err := time.Parse(pkg.DateFormatTemplateDD_MM_YYYY, search)
		if err == nil {
			rows, err = Db.Query(getTasksLimitWhereDateCommand, sql.Named("date", searchDate.Format(pkg.DateFormatTemplateYYYYMMDD)),
				sql.Named("limit", limit))
		} else {
			searchPattern := "%" + search + "%"
			rows, err = Db.Query(getTasksLimitWhereTitleOrCommentCommand,
				sql.Named("search", searchPattern), sql.Named("search", searchPattern),
				sql.Named("limit", limit))
		}
	} else {
		rows, err = Db.Query(getTasksLimitBaseCommand, sql.Named("limit", limit))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return []*Task{}, nil
	}
	return tasks, nil
}
