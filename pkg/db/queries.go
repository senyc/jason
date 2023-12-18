package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/senyc/jason/pkg/types"
)

func (db *DB) AddNewTask(newTask types.NewTask, userId string) error {
	query := "INSERT INTO tasks (user_id, title, body,  priority) VALUES (?, ?, ?, ?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, newTask.Title, newTask.Body, newTask.Priority)
	return err
}

func (db *DB) GetTaskById(userId string, taskId string) (types.Task, error) {
	var task types.Task
	var res taskRow

	query := "SELECT id, title, body, due, priority, completed FROM tasks WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return task, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId, taskId).Scan(&res.id, &res.title, &res.body, &res.due, &res.priority, &res.completed)
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

	query := "SELECT id, title, body, due, priority, completed FROM tasks WHERE user_id = ?"
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
		err := rows.Scan(&row.id, &row.title, &row.body, &row.due, &row.priority, &row.completed)
		if err != nil {
			return tasks, nil
		}
		tasks = append(tasks, removeEmptyValues(row))
	}
	return tasks, nil
}

func (db *DB) AddNewUser(newUser types.NewUser) error {
	query := "INSERT INTO users (first_name, last_name, email, account_type) VALUES (?, ?, ?, ?)"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.FirstName, newUser.LastName, newUser.Email, newUser.AccountType)
	return err
}

func (db *DB) MarkTaskCompleted(userId string, taskId string) error {
	query := "UPDATE tasks SET completed = 1 WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, taskId)
	return err
}

func (db *DB) MarkTaskIncomplete(userId string, taskId string) error {
	query := "UPDATE tasks SET completed = 0 WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, taskId)
	return err
}
