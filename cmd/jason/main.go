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
		fmt.Println(err)
	}
	server := new(server.Server)
	err = server.Start()
	defer server.Shutdown()

	if err != nil {
		panic(err)
	}
}
