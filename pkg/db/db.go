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
