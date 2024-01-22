package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/senyc/jason/pkg/types"
)

var (
	NoTasksFoundError                = errors.New("No tasks found")
	NewUserUniquenessConstraintError = errors.New("There is already an account with this email, please use another or login")
)

const uniqeConstraintErrorId = 1062

func (db *DB) GetAddedTasksCount(userId string) (int, error) {
	var addedTasksCount int
	query := "SELECT added_tasks FROM users WHERE id = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return addedTasksCount, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&addedTasksCount)
	if err == sql.ErrNoRows {
		return addedTasksCount, NoTasksFoundError
	}
	return addedTasksCount, err
}

func (db *DB) AddNewTask(newTask types.NewTaskPayload, userId string) error {
	query := "INSERT INTO tasks (user_id, id, title, body, priority, due) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Gets monotonically increasing number of tasks that have been added for the user
	taskId, err := db.GetAddedTasksCount(userId)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId, taskId, newTask.Title, newTask.Body, newTask.Priority, newTask.Due)
	return err
}

func (db *DB) GetTaskById(userId string, taskId string) (types.SqlTasksRow, error) {
	var task types.SqlTasksRow

	query := "SELECT id, title, body, due, time_created, priority, completed, completed_date FROM tasks WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return task, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId, taskId).Scan(&task.Id, &task.Title, &task.Body, &task.Due, &task.TimeCreated, &task.Priority, &task.Completed, &task.CompletedDate)
	if err == sql.ErrNoRows {
		return task, NoTasksFoundError
	}
	return task, err
}

func (db *DB) GetAllTasksByUser(uuid string) ([]types.SqlTasksRow, error) {
	var tasks []types.SqlTasksRow

	query := `SELECT id, title, body, due, time_created, priority, completed, completed_date 
	FROM tasks 
	WHERE user_id = ? 
	ORDER BY 
		CASE 
			WHEN priority BETWEEN 1 AND 5 THEN priority
			ELSE 6
		END,
	due, priority ASC`

	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return tasks, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid)
	if err != nil {
		// Handle empty row path
		return tasks, err
	}

	for rows.Next() {
		var row types.SqlTasksRow
		err = rows.Scan(&row.Id, &row.Title, &row.Body, &row.Due, &row.TimeCreated, &row.Priority, &row.Completed, &row.CompletedDate)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, row)
	}
	return tasks, nil
}

func (db *DB) GetCompletedTasks(uuid string) ([]types.SqlTasksRow, error) {
	var tasks []types.SqlTasksRow

	query := `SELECT id, title, body, due, time_created, priority, completed, completed_date 
	FROM tasks 
	WHERE user_id = ? AND completed = true 
	ORDER BY 
		CASE 
			WHEN priority BETWEEN 1 AND 5 THEN priority
			ELSE 6
		END,
	due, priority ASC`

	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return tasks, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid)
	if err != nil {
		// Handle empty row path
		return tasks, err
	}

	for rows.Next() {
		var row types.SqlTasksRow
		err = rows.Scan(&row.Id, &row.Title, &row.Body, &row.Due, &row.TimeCreated, &row.Priority, &row.Completed, &row.CompletedDate)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, row)
	}
	return tasks, nil
}
func (db *DB) GetIncompleteTasks(uuid string) ([]types.SqlTasksRow, error) {
	var tasks []types.SqlTasksRow

	query := `
	SELECT id, title, body, due, time_created, priority, completed, completed_date 
	FROM tasks 
	WHERE user_id = ? AND completed = false 
	ORDER BY 
		CASE 
			WHEN priority BETWEEN 1 AND 5 THEN priority
			ELSE 6
		END,
	due, priority ASC`

	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return tasks, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid)
	if err != nil {
		// Handle empty row path
		return tasks, err
	}

	for rows.Next() {
		var row types.SqlTasksRow
		err = rows.Scan(&row.Id, &row.Title, &row.Body, &row.Due, &row.TimeCreated, &row.Priority, &row.Completed, &row.CompletedDate)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, row)
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
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == uniqeConstraintErrorId {
				return NewUserUniquenessConstraintError
			}
		}
	}

	return err
}

func (db *DB) MarkTaskCompleted(userId string, taskId string) error {
	query := "UPDATE tasks SET completed = 1, completed_date = CURTIME() WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userId, taskId)
	if err != nil {
		return err
	}

	if v, _ := result.RowsAffected(); v == 0 {
		return NoTasksFoundError
	}
	return nil
}

func (db *DB) MarkTaskIncomplete(userId string, taskId string) error {
	query := "UPDATE tasks SET completed = 0, completed_date = NULL WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userId, taskId)
	if err != nil {
		return err
	}

	if v, _ := result.RowsAffected(); v == 0 {
		return NoTasksFoundError
	}
	return nil
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

func (db *DB) DeleteTask(userId string, taskId string) error {
	query := "DELETE FROM tasks WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userId, taskId)
	if err != nil {
		return err
	}

	if v, _ := result.RowsAffected(); v == 0 {
		return NoTasksFoundError
	}
	return nil
}

func (db *DB) EditTask(userId string, taskPayload types.EditTaskPayload) error {
	var payloads []any
	query := "UPDATE tasks SET"

	if taskPayload.Title != "" {
		query += fmt.Sprintf(" title = ?,")
		payloads = append(payloads, taskPayload.Title)
	}

	if taskPayload.Body != "" {
		query += fmt.Sprintf(" body = ?,")
		payloads = append(payloads, taskPayload.Body)
	}

	if taskPayload.Priority != 0 {
		query += fmt.Sprintf(" priority = ?,")
		payloads = append(payloads, taskPayload.Priority)
	}

	if taskPayload.Due != nil {
		query += fmt.Sprint(" due = ?,")
		payloads = append(payloads, &taskPayload.Due)
	}

	// TODO: perform this check in the handler level, we should be able to always remove last comma
	// Removes last comma otherwise sql invalid
	if query[len(query)-1] == ',' {
		query = query[:len(query)-1]
	}

	query += " WHERE user_id = ? AND id = ?"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Unpacks all of the values that are required in the built query
	_, err = stmt.Exec(append(payloads, userId, taskPayload.Id)...)
	return err
}

func (db *DB) GetEmail(userId string) (string, error) {
	var result string
	query := "SELECT email from users WHERE id = ?"

	stmt, err := db.conn.Prepare(query)

	if err != nil {
		return result, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&result)

	return result, err
}
