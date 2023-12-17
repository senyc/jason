package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() {
	pass := os.Getenv("PASSWORD")
	if pass == "" {
		panic("No password found")
	}

	dbPath := fmt.Sprintf("root:%s@tcp(localhost:3306)/jason", pass)
	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		fmt.Println("error connecting to sql server")
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("Select * from tasks where id = 1")
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(res)
}
