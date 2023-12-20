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
	// Handle empty row path
	return task, err
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
		// Handle empty row path
		return tasks, err
	}

	for rows.Next() {
		var row taskRow
		err := rows.Scan(&row.id, &row.title, &row.body, &row.due, &row.priority, &row.completed)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, removeEmptyValues(row))
	}
	return tasks, nil
}

func (db *DB) AddNewUser(newUser types.User, apiKey string) error {
	query := "INSERT INTO users (first_name, last_name, password, email, account_type, api_key) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.FirstName, newUser.LastName, newUser.Password, newUser.Email, newUser.AccountType, apiKey)
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

func (db *DB) GetApiKey(userAuth types.UserLogin) (string, error) {
	var apiKey string
	query := "SELECT api_key FROM users WHERE email = ? AND password = ?"

	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return apiKey, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userAuth.Email, userAuth.Password).Scan(&apiKey)
	return apiKey, err
}

func (db *DB) GetUserIdFromApiKey(apiKey string) (string, error) {
	var userId string

	query := "SELECT id FROM users WHERE api_key = ?"

	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return userId, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(apiKey).Scan(&userId)
	return userId, err
}

func (db *DB) GetPasswordFromLogin(login string) (string, error) {
	var result string
	query := "SELECT password from users WHERE email = ?"

	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return result, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(login).Scan(&result)

	return result, err
}
