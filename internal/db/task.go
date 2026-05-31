// Package db
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antonlearn/go-final-project/pkg/config"
	"github.com/antonlearn/go-final-project/pkg/format"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	result, err := config.Config.DBConnect.Exec(`INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (:date, :title, :comment, :repeat)`, sql.Named("date", task.Date),
		sql.Named("title", task.Title), sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	config.Config.Logger.Printf("Task %v was added successfully\n", task)
	return result.LastInsertId()
}

func GetTasks(search string) ([]*Task, error) {
	var (
		tasks      []*Task
		rows       *sql.Rows
		err        error
		searchDate time.Time
	)
	search = strings.TrimSpace(search)
	if search != "" {
		searchDate, err = time.Parse(format.DateFormatTemplateDDMMYYYY, search)
		if err == nil {
			rows, err = config.Config.DBConnect.Query(`SELECT id, date, title, comment, repeat 
			FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit`,
				sql.Named("date", searchDate.Format(format.DateFormatTemplateYYYYMMDD)),
				sql.Named("limit", config.Config.MaxNumTasks))
		} else {
			rows, err = config.Config.DBConnect.Query(`SELECT id, date, title, comment, repeat 
			FROM scheduler WHERE LOWER(title) LIKE CONCAT('%', LOWER(:search), '%') 
			OR LOWER(comment) LIKE CONCAT('%', LOWER(:search), '%') ORDER BY date 
			LIMIT :limit`, sql.Named("search", search), sql.Named("search", search),
				sql.Named("limit", config.Config.MaxNumTasks))
		}
	} else {
		rows, err = config.Config.DBConnect.Query(`SELECT id, date, title, comment, repeat 
			FROM scheduler ORDER BY date LIMIT :limit`,
			sql.Named("limit", config.Config.MaxNumTasks))
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
	config.Config.Logger.Printf("List of upcoming tasks %v has been successfully created\n", tasks)
	return tasks, nil
}

func GetTask(id int) (*Task, error) {
	var task Task
	err := config.Config.DBConnect.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler 
		WHERE id = :id`, sql.Named("id", id)).Scan(&task.ID, &task.Date,
		&task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("request parameters: no task with ID %d", id)
		}
		return nil, err
	}
	config.Config.Logger.Printf("Task %v with ID %d received successfully\n", task, id)
	return &task, nil
}

func UpdateTask(task *Task) error {
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}
	result, err := config.Config.DBConnect.Exec(`UPDATE scheduler SET date = :date, title = :title, 
		comment = :comment, repeat = :repeat WHERE id = :id`,
		sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat),
		sql.Named("id", id))
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
		config.Config.Logger.Printf("Task %v with ID %d successfully updated\n", task, id)
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}

func DeleteTask(id int) error {
	_, err := config.Config.DBConnect.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", id))
	if err != nil {
		return err
	}
	config.Config.Logger.Printf("Task with ID %d successfully deleted\n", id)
	return nil
}

func UpdateDateTask(nextDate string, id int) error {
	result, err := config.Config.DBConnect.Exec(`UPDATE scheduler SET date = :date WHERE id = :id`,
		sql.Named("date", nextDate), sql.Named("id", id))
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
		config.Config.Logger.Printf("Task date with ID %d successfully updated\n", id)
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}
