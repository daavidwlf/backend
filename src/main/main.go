package main

import (
	"backend/src/db"
	"backend/src/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port_env := os.Getenv("BACKEND_PORT")

	port := server.CreateServer(":" + port_env)
	db.ConnectDB()

	server.Run(port)
}
