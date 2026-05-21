package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
	getTaskByIdCommand    = `SELECT * FROM scheduler WHERE id = :id`
	updateTaskCommand     = `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	deleteTaskCommand     = `DELETE FROM scheduler WHERE id = :id`
	updateDateTaskCommand = `UPDATE scheduler SET date = :date WHERE id = :id`
)

func AddTask(task *Task) (int64, error) {
	result, err := Db.Exec(addTaskCommand, sql.Named("date", task.Date), sql.Named("title",
		task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetTasks(search string) ([]*Task, error) {
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
				sql.Named("limit", pkg.MaxNumTasks))
		} else {
			searchPattern := "%" + search + "%"
			rows, err = Db.Query(getTasksLimitWhereTitleOrCommentCommand,
				sql.Named("search", searchPattern), sql.Named("search", searchPattern),
				sql.Named("limit", pkg.MaxNumTasks))
		}
	} else {
		rows, err = Db.Query(getTasksLimitBaseCommand, sql.Named("limit", pkg.MaxNumTasks))
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

func GetTask(id int) (*Task, error) {
	var task Task
	err := Db.QueryRow(getTaskByIdCommand, sql.Named("id", id)).Scan(&task.ID, &task.Date,
		&task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("request parameters: no task with ID %d", id)
		}
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}
	result, err := Db.Exec(updateTaskCommand, sql.Named("date", task.Date),
		sql.Named("title", task.Title), sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat), sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	switch count {
	case 0:
		return fmt.Errorf("request parameters: incorrect ID %d for updating task", id)
	case 1:
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}

func DeleteTask(id int) error {
	_, err := Db.Exec(deleteTaskCommand, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

func UpdateDateTask(nextDate string, id int) error {
	result, err := Db.Exec(updateDateTaskCommand, sql.Named("date", nextDate), sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	switch count {
	case 0:
		return fmt.Errorf("request parameters: incorrect ID %d for updating task", id)
	case 1:
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}
