package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

type Task struct {
	id       string
	title    string
	body     sql.NullString
	due      sql.NullString
	priority sql.NullString
}

func (db *DB) GetTaskById(userId string, taskId string) (Task, error) {
	var task Task

	query := "SELECT id, title, body, due, priority FROM tasks WHERE user_id = ? AND id = ? "
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return task, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId, taskId).Scan(&task.id, &task.title, &task.body, &task.due, &task.priority)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nothing returned")
		}
		return task, err
	}
	return task, nil
}

func (db *DB) GetAllTasksByUser(userId string) (Task, error) {
	var task Task

	query := "SELECT id, title, body, due, priority FROM tasks WHERE id = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return task, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&task.id, &task.title, &task.body, &task.due, &task.priority)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nothing returned")
		}
		return task, err
	}
	return task, nil
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
