package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/senyc/jason/pkg/types"
)

type DB struct {
	conn *sql.DB
}

type taskRow struct {
	id       string
	title    string
	body     sql.NullString
	due      sql.NullString
	priority sql.NullString
}

func removeEmptyValues(task taskRow) types.Task {
	result := types.Task{
		Id:       task.id,
		Title:    task.title,
		Body:     "",
		Due:      "",
		Priority: "",
	}

	if task.body.Valid {
		result.Body = task.body.String
	}

	if task.due.Valid {
		result.Due = task.due.String
	}

	if task.priority.Valid {
		result.Priority = task.priority.String
	}

	return result
}

func (db *DB) GetTaskById(userId string, taskId string) (types.Task, error) {
	var task types.Task
	var res taskRow

	query := "SELECT id, title, body, due, priority FROM tasks WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return task, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId, taskId).Scan(&res.id, &res.title, &res.body, &res.due, &res.priority)
	task = removeEmptyValues(res)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nothing returned")
		}
		return task, err
	}
	return task, nil
}

func (db *DB) GetAllTasksByUser(userId string) ([]types.Task, error) {
	var tasks []types.Task

	query := "SELECT id, title, body, due, priority FROM tasks WHERE user_id = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return tasks, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nothing returned")
		}
		return tasks, err
	}

	for rows.Next() {
		var row taskRow
		err := rows.Scan(&row.id, &row.title, &row.body, &row.due, &row.priority)
		if err != nil {
			return tasks, nil
		}
		tasks = append(tasks, removeEmptyValues(row))
	}
	return tasks, nil
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

	dbPath := fmt.Sprintf("%s:%s@tcp(%s:%s)/jason", user, pass, domain, port)

	connection, err := sql.Open("mysql", dbPath)
	if err != nil {
		return err
	}

	db.conn = connection
	return db.conn.Ping()
}
