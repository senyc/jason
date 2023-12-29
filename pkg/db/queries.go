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

func (db *DB) GetCompletedTasks(uuid string) ([]types.CompletedTask, error) {
	var tasks []types.CompletedTask

	query := "SELECT id, title, body, due, priority, completed_date FROM tasks WHERE user_id = ? AND completed = true"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return tasks, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nothing returned")
		}
		// Handle empty row path
		return tasks, err
	}

	for rows.Next() {
		var row completedTaskRow
		err := rows.Scan(&row.id, &row.title, &row.body, &row.due, &row.priority)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, removeEmptyValuesComplete(row))
	}
	return tasks, nil
}
func (db *DB) GetIncompleteTasks(uuid string) ([]types.IncompleteTask, error) {
	var tasks []types.IncompleteTask

	query := "SELECT id, title, body, due, priority, FROM tasks WHERE user_id = ? AND completed = false"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return tasks, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nothing returned")
		}
		// Handle empty row path
		return tasks, err
	}

	for rows.Next() {
		var row incompleteTaskRow
		err := rows.Scan(&row.id, &row.title, &row.body, &row.due, &row.priority)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, removeEmptyValuesIncomplete(row))
	}
	return tasks, nil
}

func (db *DB) AddNewUser(newUser types.User) error {
	query := "INSERT INTO users (password, email, account_type) VALUES (?, ?, ?)"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.Password, newUser.Email, newUser.AccountType)
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

func (db *DB) AddApiKey(userLogin string, apiKey string) error {
	query := "UPDATE users SET api_key = ? WHERE email = ?"

	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(apiKey, userLogin)
	return err
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

func (db *DB) GetUuidFromEmail(email string) (string, error) {
	var result string
	query := "SELECT id from users WHERE email = ?"

	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return result, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(email).Scan(&result)

	return result, err
}
