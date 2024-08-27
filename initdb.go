package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// initDB creates the required tables if they don't exist
func initDB() {
	fmt.Println("Server: Creating tables...")

	// Create test table
	testTableQuery := `CREATE TABLE IF NOT EXISTS test (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	)`

	_, err := db.Exec(testTableQuery)
	if err != nil {
		log.Fatal("Server: Error creating test table: ", err.Error())
	}

	// Insert the first row into the test table
	insertTestRowQuery := `INSERT INTO test (id, name) VALUES (1, 'test')`
	_, err = db.Exec(insertTestRowQuery)
	if err != nil {
		log.Fatal("Server: Error inserting first row into test table: ", err.Error())
	}

	// Create users table
	usersTableQuery := `CREATE TABLE IF NOT EXISTS users (
		UserID INT AUTO_INCREMENT PRIMARY KEY,
		FirstName VARCHAR(255) NOT NULL,
		LastName VARCHAR(255) NOT NULL,
		Email VARCHAR(255) NOT NULL UNIQUE,
		Password VARCHAR(255) NOT NULL,
		Created DATE NOT NULL
	)`

	_, err = db.Exec(usersTableQuery)
	if err != nil {
		log.Fatal("Server: Error creating users table: ", err.Error())
	}

	fmt.Println("Server: Tables created and initial row inserted successfully")
}
