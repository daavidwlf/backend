package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func connectDB() {

	//dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(db_mysql:%s)/%s", dbUser, dbPass, dbPort, dbName)

	fmt.Println(dsn)
	fmt.Println("Server: Opening database")

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Server: Couldn't open database: ", err.Error())
		return
	}

	fmt.Println("Server: Database opend")

	//test if connection to db was established
	for i := 0; i < 5; i++ {
		err = db.Ping()

		if err == nil {

			query := "SELECT id, name FROM test where id = ?"
			var id int
			var name string

			// perform a a test query - needs to be removed later
			err = db.QueryRow(query, 1).Scan(&id, &name)

			if err != nil {
				log.Fatal("Server: Error while perfoming Query: ", err.Error())
			}

			fmt.Println("Server: Sucessfully performed Query")
			fmt.Printf("Results: id: %d , name: %s", id, name)

			break
		}

		//retrying after 5 seconds
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Server: Unable to ping database: ", err.Error())
	}

	fmt.Println("Server: Succesfully connected to Database")
}

func registerUser(usr user) error {

	var mail string

	err := db.QueryRow(`SELECT Email FROM users where Email = ?`, usr.Email).Scan(&mail)

	if err == nil {
		return errors.New("user already exists")
	}

	if err != sql.ErrNoRows {
		return errors.New("couldn't execute user search in database: " + err.Error())
	}

	//create props for new user
	var newID uuid.UUID
	var IDerr error
	usr.ID, IDerr = uuid.NewUUID()

	if IDerr != nil {
		return errors.New("couldn't generate UUID: " + IDerr.Error())
	}

	usr.Created = int(time.Now().Unix())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("couldn't hash password: " + err.Error())
	}

	usr.Password = string(hashedPassword)

	var rows *sql.Rows
	rows, err = db.Query(`INSERT INTO users (UserID, FirstName, LastName, Email, Password, Created) VALUES (?, ?, ?, ?, ?, ?)`, usr.ID, usr.FirstName, usr.LastName, usr.Email, usr.Password, usr.Created)

	if err != nil {
		return errors.New("couldn't execute user creation on db: " + err.Error())
	}

	defer rows.Close()

	fmt.Println("Server: New user created: ID: ", newID)

	return err
}

func loginUser(usr user) (string, error) {

	var reqPassword string
	var usrID string

	err := db.QueryRow(`SELECT UserID, Password FROM users where email = ?`, usr.Email).Scan(&usrID, &reqPassword)

	if err == sql.ErrNoRows {
		return "", errors.New("email doesn't exist")
	}

	if err != nil {
		return "", errors.New("error while logging in" + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(reqPassword), []byte(usr.Password))

	return usrID, nil
}

func getUserByID(usrID string) (*user, error) {

	var usr user

	err := db.QueryRow(`SELECT UserID, FirstName, LastName, Email, Created FROM users WHERE UserID = ?`, usrID).Scan(&usr.ID, &usr.FirstName, &usr.LastName, &usr.Email, &usr.Created)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, errors.New("error occured getting user from db" + err.Error())
	}

	return &usr, nil
}
