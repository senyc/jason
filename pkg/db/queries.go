package db

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/senyc/jason/pkg/types"
)

var (
	NoTasksFoundError                = errors.New("No tasks found")
	NewUserUniquenessConstraintError = errors.New("There is already an account with this email, please use another or login")
)

const uniqeConstraintErrorId = 1062

func (db *DB) AddNewTask(newTask types.NewTaskPayload, userId string) error {
	query := "INSERT INTO tasks (user_id, title, body, priority, due) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, newTask.Title, newTask.Body, newTask.Priority, newTask.Due)
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

	query := "SELECT id, title, body, due, time_created, priority, completed, completed_date FROM tasks WHERE user_id = ? ORDER BY due ASC"
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

	query := "SELECT id, title, body, due, time_created, priority, completed, completed_date FROM tasks WHERE user_id = ? AND completed = true ORDER BY due ASC"
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

	query := "SELECT id, title, body, due, time_created, priority, completed, completed_date FROM tasks WHERE user_id = ? AND completed = false ORDER BY due ASC"
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
