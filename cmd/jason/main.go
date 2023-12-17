package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/senyc/jason/pkg/server"
)

func main() {
	// Loads environment variables 
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error(), "No environment variables found")
	}

	server.Start()
}
