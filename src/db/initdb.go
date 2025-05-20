package db

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// initDB creates the required tables if they don't exist
func initDB() {
	fmt.Println("Server: Creating tables...")

	// Test table
	testTableQuery := `CREATE TABLE IF NOT EXISTS test (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	)`

	_, err := db.Exec(testTableQuery)
	if err != nil {
		log.Fatal("Server: Error creating test table: ", err.Error())
	}

	insertTestRowQuery := `INSERT INTO test (id, name) VALUES (1, 'test')`
	_, err = db.Exec(insertTestRowQuery)
	if err != nil {
		log.Fatal("Server: Error inserting first row into test table: ", err.Error())
	}

	// Create admins table
	adminsTableQuery := `CREATE TABLE IF NOT EXISTS admins (
		AdminID varchar(36) NOT NULL PRIMARY KEY,
		Email text NOT NULL,
		UserName text NOT NULL,
		Password text NOT NULL,
		Created int NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`

	_, err = db.Exec(adminsTableQuery)
	if err != nil {
		log.Fatal("Server: Error creating admins table: ", err.Error())
	}

	// Insert initial rows into the admins table
	insertAdminsQuery := `INSERT INTO admins (AdminID, Email, UserName, Password, Created) VALUES
	('99278b45-63d3-11ef-9353-0242c0a8b502', 'julian.boehne@web.de', 'Julian', '$2a$12$vQmM9YShnUlX9ZZFRXwNOuRkbNmi8dSMjHfx0wKekXJZeoeGT4dvO', 1724694578),
	('d23d9df9-63d3-11ef-9353-0242c0a8b502', 'wolf_david@gmx.de', 'David', '$2a$12$foK/kJYQn6QjlOTFXIw9FODo2motgflWuM2xTA0agV/HqiVQ2inCu', 1724694649);`

	_, err = db.Exec(insertAdminsQuery)
	if err != nil {
		log.Fatal("Server: Error inserting initial rows into admins table: ", err.Error())
	}

	// Create users table with UserID as PRIMARY KEY and UNIQUE
	usersTableQuery := `CREATE TABLE IF NOT EXISTS users (
		UserID varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL PRIMARY KEY,
		FirstName text NOT NULL,
		LastName text NOT NULL,
		Email text NOT NULL,
		Password text NOT NULL,
		Created int NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`

	_, err = db.Exec(usersTableQuery)
	if err != nil {
		log.Fatal("Server: Error creating users table: ", err.Error())
	}

	fmt.Println("Server: Tables created and initial data inserted successfully")
}
