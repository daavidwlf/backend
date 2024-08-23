package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("BACKEND_PORT")

	server := createServer(":" + port)
	connectDB()
	server.run()
}
