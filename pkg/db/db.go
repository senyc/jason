package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/senyc/jason/pkg/types"
)

type DB struct {
	conn *sql.DB
}

type taskRow struct {
	id        int
	title     string
	body      sql.NullString
	due       sql.NullTime
	priority  sql.NullInt16
	completed bool
}

func removeEmptyValues(task taskRow) types.Task {
	result := types.Task{
		Id:        task.id,
		Title:     task.title,
		Body:      "",
		Due:       time.Time{},
		Priority:  3,
		Completed: task.completed,
	}

	if task.body.Valid {
		result.Body = task.body.String
	}

	if task.due.Valid {
		result.Due = task.due.Time
	}

	if task.priority.Valid {
		result.Priority = task.priority.Int16
	}

	return result
}

type completedTaskRow struct {
	id            int
	title         string
	body          sql.NullString
	due           sql.NullTime
	priority      sql.NullInt16
	completedDate sql.NullString
}

type incompleteTaskRow struct {
	id       int
	title    string
	body     sql.NullString
	due      sql.NullTime
	priority sql.NullInt16
}

func removeEmptyValuesIncomplete(task incompleteTaskRow) types.IncompleteTask {
	result := types.IncompleteTask{
		Id:       task.id,
		Title:    task.title,
		Body:     "",
		Due:      time.Time{},
		Priority: 3,
	}

	if task.body.Valid {
		result.Body = task.body.String
	}

	if task.due.Valid {
		result.Due = task.due.Time
	}

	if task.priority.Valid {
		result.Priority = task.priority.Int16
	}
	return result
}

// TODO: need to actually add the completed date functionality
func removeEmptyValuesComplete(task completedTaskRow) types.CompletedTask {
	result := types.CompletedTask{
		Id:            task.id,
		Title:         task.title,
		Body:          "",
		Due:           time.Time{},
		Priority:      3,
		CompletedDate: "",
	}

	if task.body.Valid {
		result.Body = task.body.String
	}

	if task.due.Valid {
		result.Due = task.due.Time
	}

	if task.priority.Valid {
		result.Priority = task.priority.Int16
	}
	return result
}

func (db *DB) Connect() error {
	pass := os.Getenv("DB_PASS")
	user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	domain := os.Getenv("DOMAIN")

	if pass == "" {
		return errors.New("No password found")
	}

	if user == "" {
		user = "root"
	}

	if port == "" {
		port = "3306"
	}

	if domain == "" {
		domain = "localhost"
	}

	dbPath := fmt.Sprintf("%s:%s@tcp(%s:%s)/jason?parseTime=true", user, pass, domain, port)

	connection, err := sql.Open("mysql", dbPath)
	if err != nil {
		return err
	}

	db.conn = connection
	return db.conn.Ping()
}
