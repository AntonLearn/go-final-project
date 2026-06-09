// Package db implements data persistence layers and CRUD query logic
// for scheduling entries utilizing parameterized SQLite statements.
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antonlearn/go-final-project/internal/model"
	"github.com/antonlearn/go-final-project/pkg/format"
)

// AddTask inserts a new task record into the scheduler table and returns its auto-incremented ID.
func (s *Store) AddTask(task *model.Task) (int64, error) {
	result, err := s.db.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		s.logger.Errorf("Failed to add task: %v", err)
		return 0, err
	}

	id, _ := result.LastInsertId()
	s.logger.Infof("Task added successfully with ID %d", id)
	return id, nil
}

// GetTasks retrieves a filtered collection of tasks up to the configured limit constraint.
// It matches exact dates if the query conforms to DD.MM.YYYY, performs a case-insensitive
// full-text search on titles and comments otherwise, or returns all tasks if empty.
func (s *Store) GetTasks(search string) ([]*model.Task, error) {
	var (
		tasks      []*model.Task
		rows       *sql.Rows
		err        error
		searchDate time.Time
	)

	search = strings.TrimSpace(search)

	if search != "" {
		// Attempt to parse the search query as an explicit calendar date constraint.
		searchDate, err = time.Parse(format.DDMMYYYY, search)
		if err == nil {
			rows, err = s.db.Query(`
				SELECT id, date, title, comment, repeat 
				FROM scheduler 
				WHERE date = :date 
				ORDER BY date 
				LIMIT :limit`,
				sql.Named("date", searchDate.Format(format.YYYYMMDD)),
				sql.Named("limit", s.maxNumTasks),
			)
		} else {
			// Fall back to a fuzzy lookup across text fields if date parsing fails.
			rows, err = s.db.Query(`
				SELECT id, date, title, comment, repeat 
				FROM scheduler 
				WHERE LOWER(title) LIKE '%' || LOWER(:search) || '%' 
				   OR LOWER(comment) LIKE '%' || LOWER(:search) || '%' 
				ORDER BY date 
				LIMIT :limit`,
				sql.Named("search", search),
				sql.Named("limit", s.maxNumTasks),
			)
		}
	} else {
		// Fetch un-filtered records ordered sequentially by due date.
		rows, err = s.db.Query(`
			SELECT id, date, title, comment, repeat 
			FROM scheduler 
			ORDER BY date 
			LIMIT :limit`,
			sql.Named("limit", s.maxNumTasks),
		)
	}

	if err != nil {
		s.logger.Errorf("Failed to query tasks with search='%s': %v", search, err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows result set and map columns to model entities.
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			s.logger.Errorf("Failed to scan task row: %v", err)
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	// Catch any internal errors encountered during row stream iteration.
	if err = rows.Err(); err != nil {
		s.logger.Errorf("Error iterating over tasks: %v", err)
		return nil, err
	}

	s.logger.Infof("Retrieved %d tasks (search: '%s')", len(tasks), search)
	return tasks, nil
}

// GetTask fetches a single task record uniquely identified by its primary key.
func (s *Store) GetTask(id int) (*model.Task, error) {
	var task model.Task
	err := s.db.QueryRow(`
		SELECT id, date, title, comment, repeat 
		FROM scheduler WHERE id = :id`, sql.Named("id", id)).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("request parameters: no task with ID %d", id)
		}
		s.logger.Errorf("Failed to get task %d: %v", id, err)
		return nil, err
	}

	s.logger.Infof("Task %d retrieved successfully", id)
	return &task, nil
}

// UpdateTask modifies all properties of an existing row targeting a specified integer ID.
func (s *Store) UpdateTask(task *model.Task) error {
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(`
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
		s.logger.Errorf("Failed to update task %d: %v", id, err)
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Validate exactly one matching row was affected by the execution command.
	switch count {
	case 0:
		return fmt.Errorf("request parameters: incorrect ID %d for updating task", id)
	case 1:
		s.logger.Infof("Task %d successfully updated", id)
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}

// DeleteTask permanently drops a task row from the database using its primary key.
func (s *Store) DeleteTask(id int) error {
	_, err := s.db.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", id))
	if err != nil {
		s.logger.Errorf("Failed to delete task %d: %v", id, err)
		return err
	}

	s.logger.Infof("Task %d successfully deleted", id)
	return nil
}

// UpdateDateTask alters exclusively the execution date column of a recurring task instance.
func (s *Store) UpdateDateTask(nextDate string, id int) error {
	result, err := s.db.Exec(`
		UPDATE scheduler SET date = :date WHERE id = :id`,
		sql.Named("date", nextDate),
		sql.Named("id", id),
	)
	if err != nil {
		s.logger.Errorf("Failed to update date for task %d: %v", id, err)
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Ensure structural integrity by confirming that only a single targeted record changed.
	switch count {
	case 0:
		return fmt.Errorf("request parameters: incorrect ID %d for updating task", id)
	case 1:
		s.logger.Infof("Date for task %d successfully updated to %s", id, nextDate)
		return nil
	default:
		return errors.New("incorrect unexpected task update")
	}
}
