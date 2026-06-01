// Package db provides database operations for the task scheduler
// using SQLite as the storage backend.
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
	"github.com/antonlearn/go-final-project/pkg/logger"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask inserts a new task into the database and returns its ID.
func AddTask(task *Task) (int64, error) {
	result, err := config.Config.DBConnect.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		logger.Errorf("Failed to add task: %v", err)
		return 0, err
	}

	id, _ := result.LastInsertId()
	logger.Infof("Task added successfully with ID %d", id)
	return id, nil
}

// GetTasks retrieves tasks with optional search filter.
func GetTasks(search string) ([]*Task, error) {
	var (
		tasks      []*Task
		rows       *sql.Rows
		err        error
		searchDate time.Time
	)

	search = strings.TrimSpace(search)

	if search != "" {
		// Try to parse as date first
		searchDate, err = time.Parse(format.DateFormatTemplateDDMMYYYY, search)
		if err == nil {
			rows, err = config.Config.DBConnect.Query(`
				SELECT id, date, title, comment, repeat 
				FROM scheduler 
				WHERE date = :date 
				ORDER BY date 
				LIMIT :limit`,
				sql.Named("date", searchDate.Format(format.DateFormatTemplateYYYYMMDD)),
				sql.Named("limit", config.Config.MaxNumTasks),
			)
		} else {
			// Search by title or comment
			rows, err = config.Config.DBConnect.Query(`
				SELECT id, date, title, comment, repeat 
				FROM scheduler 
				WHERE LOWER(title) LIKE '%' || LOWER(:search) || '%' 
				   OR LOWER(comment) LIKE '%' || LOWER(:search) || '%' 
				ORDER BY date 
				LIMIT :limit`,
				sql.Named("search", search),
				sql.Named("limit", config.Config.MaxNumTasks),
			)
		}
	} else {
		// Return all tasks
		rows, err = config.Config.DBConnect.Query(`
			SELECT id, date, title, comment, repeat 
			FROM scheduler 
			ORDER BY date 
			LIMIT :limit`,
			sql.Named("limit", config.Config.MaxNumTasks),
		)
	}

	if err != nil {
		logger.Errorf("Failed to query tasks with search='%s': %v", search, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			logger.Errorf("Failed to scan task row: %v", err)
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		logger.Errorf("Error iterating over tasks: %v", err)
		return nil, err
	}

	logger.Infof("Retrieved %d tasks (search: '%s')", len(tasks), search)
	return tasks, nil
}

// GetTask retrieves a single task by ID.
func GetTask(id int) (*Task, error) {
	var task Task
	err := config.Config.DBConnect.QueryRow(`
		SELECT id, date, title, comment, repeat 
		FROM scheduler WHERE id = :id`, sql.Named("id", id)).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("request parameters: no task with ID %d", id)
		}
		logger.Errorf("Failed to get task %d: %v", id, err)
		return nil, err
	}

	logger.Infof("Task %d retrieved successfully", id)
	return &task, nil
}

// UpdateTask updates an existing task.
func UpdateTask(task *Task) error {
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}

	result, err := config.Config.DBConnect.Exec(`
		UPDATE scheduler 
		SET date = :date, title = :title, comment = :comment, repeat = :repeat 
		WHERE id = :id`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", id),
	)
	if err != nil {
		logger.Errorf("Failed to update task %d: %v", id, err)
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
		logger.Infof("Task %d successfully updated", id)
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}

// DeleteTask removes a task by ID.
func DeleteTask(id int) error {
	_, err := config.Config.DBConnect.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", id))
	if err != nil {
		logger.Errorf("Failed to delete task %d: %v", id, err)
		return err
	}

	logger.Infof("Task %d successfully deleted", id)
	return nil
}

// UpdateDateTask updates only the date of a recurring task.
func UpdateDateTask(nextDate string, id int) error {
	result, err := config.Config.DBConnect.Exec(`
		UPDATE scheduler SET date = :date WHERE id = :id`,
		sql.Named("date", nextDate),
		sql.Named("id", id),
	)
	if err != nil {
		logger.Errorf("Failed to update date for task %d: %v", id, err)
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
		logger.Infof("Date for task %d successfully updated to %s", id, nextDate)
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}
