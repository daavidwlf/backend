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

	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	fmt.Println(dsn)
	fmt.Println("Server: Opening database")

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Server: Couldn't open database: ", err.Error())
		return
	}

	fmt.Println("Server: Database opend")

	// test if connection to db was established
	for i := 0; i < 5; i++ {
		err = db.Ping()

		if err == nil {
			// check first row
			query := "SELECT id, name FROM test LIMIT 1"
			var id int
			var name string

			// perform a a test query - needs to be removed later
			err = db.QueryRow(query).Scan(&id, &name)

			if err != nil {
				fmt.Println("Server: Table not found, initializing database...")
				initDB()
			}

			err = db.QueryRow(query).Scan(&id, &name)

			if err != nil {
				log.Fatal("Server: Error while perfoming Query: ", err.Error())
			}

			fmt.Println("Server: Sucessfully performed Query")
			fmt.Printf("Results: id: %d , name: %s\n", id, name)

			break
		}

		// retrying after 5 seconds
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Server: Unable to ping database: ", err.Error())
	}

	fmt.Println("Server: Succesfully connected to Database")
}

func registerUser(usr registerUserRequest) error {

	var mail string

	err := db.QueryRow(`SELECT Email FROM users where Email = ?`, usr.Email).Scan(&mail)

	if err == nil {
		return errors.New("user already exists")
	}

	if err != sql.ErrNoRows {
		return errors.New("couldn't execute user search in database: " + err.Error())
	}

	// create new user
	var newUser user
	var IDerr error
	newUser.ID, IDerr = uuid.NewUUID()

	if IDerr != nil {
		return errors.New("couldn't generate UUID: " + IDerr.Error())
	}

	newUser.Created = int(time.Now().Unix())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("couldn't hash password: " + err.Error())
	}

	newUser.Password = string(hashedPassword)

	newUser.Email = usr.Email
	newUser.FirstName = usr.FirstName
	newUser.LastName = usr.LastName

	var rows *sql.Rows
	rows, err = db.Query(`INSERT INTO users (UserID, FirstName, LastName, Email, Password, Created) VALUES (?, ?, ?, ?, ?, ?)`, newUser.ID, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password, newUser.Created)

	if err != nil {
		return errors.New("couldn't execute user creation on db: " + err.Error())
	}

	defer rows.Close()

	fmt.Println("Server: New user created: ID: ", newUser.ID)

	return err
}

func getUserByID(usrID string) (*user, error) {

	var usr user

	err := db.QueryRow(`SELECT UserID, FirstName, LastName, Email, Created FROM users WHERE UserID = ?`, usrID).Scan(&usr.ID, &usr.FirstName, &usr.LastName, &usr.Email, &usr.Created)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, errors.New("error occured getting user from db " + err.Error())
	}

	return &usr, nil
}

func loginUser(usr loginUserRequest) (string, error) {
	query := `SELECT UserID, Password FROM users where email = ?`
	return loginHelper(usr.Email, usr.Password, query, "UserID")
}

func loginAdmin(adm loginAdminRequest) (string, error) {
	query := `SELECT AdminID, Password FROM admins where email = ?`
	return loginHelper(adm.Email, adm.Password, query, "AdminID")
}

func loginHelper(email, password, query, idField string) (string, error) {
	var requiredPassword string
	var userID string

	err := db.QueryRow(query, email).Scan(&userID, &requiredPassword)

	if err == sql.ErrNoRows {
		return "", errors.New("email doesn't exist")
	}

	if err != nil {
		return "", errors.New("error while logging in " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(requiredPassword), []byte(password))

	if err != nil {
		return "", errors.New("wrong password")
	}

	return userID, nil
}

func getAdminByID(admID string) (*admin, error) {

	var adm admin

	err := db.QueryRow(`SELECT AdminID, Email, UserName, Created FROM admins WHERE AdminID = ?`, admID).Scan(&adm.ID, &adm.Email, &adm.UserName, &adm.Created)

	if err == sql.ErrNoRows {
		return nil, errors.New("admin not found")
	}

	if err != nil {
		return nil, errors.New("error occured getting admin from db " + err.Error())
	}

	return &adm, nil
}

func getMultibleAdmins(quantity int) (*[]admin, error) {
	var adminList []admin

	rows, err := db.Query(`SELECT AdminID, Email, UserName, Created FROM admins LIMIT ?`, quantity)

	if err != nil {
		return nil, errors.New("unable to get admins " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var current admin

		err := rows.Scan(&current.ID, &current.Email, &current.UserName, &current.Created)

		if err != nil {
			return nil, errors.New("error while appending admins " + err.Error())
		}

		adminList = append(adminList, current)
	}

	return &adminList, nil

}

func editAdmin(adminID string, edit *editAdminRequest) (*editAdminRequest, error) {

	fmt.Println(adminID)

	fmt.Println(edit.Email, edit.UserName)

	result, err := db.Exec(`UPDATE admins SET Username = ?, Email = ? WHERE AdminID = ?`, edit.UserName, edit.Email, adminID)

	if err != nil {
		return nil, errors.New("error while updating db " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.New("error while checking affected rows: " + err.Error())
	}

	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}

	err = db.QueryRow(`SELECT Username, Email FROM admins WHERE AdminID = ?`, adminID).Scan(&edit.UserName, &edit.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no admin found after update with the given AdminID")
		}
		return nil, errors.New("error while fetching updated admin: " + err.Error())
	}

	return edit, nil
}
