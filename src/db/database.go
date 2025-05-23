package db

import (
	customTypes "backend/src/types"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func ConnectDB() {

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

func EditPerson(person customTypes.Person, id string, usr *customTypes.EditUserRequest, adm *customTypes.EditAdminRequest) (string, error) {

	var result sql.Result
	var err error

	switch person {
	case customTypes.USER:
		result, err = db.Exec(`UPDATE users SET FirstName = ?, LastName = ?, Email = ? WHERE UserID = ?`, usr.FirstName, usr.LastName, usr.Email, id)
	case customTypes.ADMIN:
		result, err = db.Exec(`UPDATE admins SET UserName = ?, Email = ? WHERE AdminID = ?`, adm.UserName, adm.Email, id)
	default:
		return "", errors.New("invalid person type")
	}

	if err != nil {
		return "", errors.New("error while updating db " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return "", errors.New("error while checking affected rows: " + err.Error())
	}

	if rowsAffected == 0 {
		return "", errors.New("no rows affected")
	}

	return id, nil
}

func GetMultiblePersons(person customTypes.Person, quantity int) (*[]customTypes.User, *[]customTypes.Admin, error) {

	var adminList []customTypes.Admin
	var userList []customTypes.User

	var rows *sql.Rows
	var err error

	switch person {
	case customTypes.USER:
		rows, err = db.Query(`SELECT UserID, FirstName, LastName, Email, Created FROM users LIMIT ?`, quantity)
	case customTypes.ADMIN:
		rows, err = db.Query(`SELECT AdminID, Email, UserName, Created FROM admins LIMIT ?`, quantity)
	default:
		return nil, nil, errors.New("invalid person type")
	}

	if err != nil {
		return nil, nil, errors.New("unable to perform query " + err.Error())
	}

	defer rows.Close()

	if person == customTypes.USER {
		for rows.Next() {
			var current customTypes.User

			err := rows.Scan(&current.ID, &current.FirstName, &current.LastName, &current.Email, &current.Created)

			if err != nil {
				return nil, nil, errors.New("error while appending users " + err.Error())
			}

			userList = append(userList, current)
		}

		return &userList, nil, nil
	}

	if person == customTypes.ADMIN {
		for rows.Next() {
			var current customTypes.Admin

			err := rows.Scan(&current.ID, &current.Email, &current.UserName, &current.Created)

			if err != nil {
				return nil, nil, errors.New("error while appending admins " + err.Error())
			}

			adminList = append(adminList, current)
		}

		return nil, &adminList, nil
	}

	return nil, nil, errors.New("invalid person type")

}

func RegisterUser(usr customTypes.RegisterUserRequest) error {

	var mail string

	err := db.QueryRow(`SELECT Email FROM users where Email = ?`, usr.Email).Scan(&mail)

	if err == nil {
		return errors.New("user already exists")
	}

	if err != sql.ErrNoRows {
		return errors.New("couldn't execute user search in database: " + err.Error())
	}

	// create new user
	var newUser customTypes.User
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

func GetUserByID(usrID string) (*customTypes.User, error) {

	var usr customTypes.User

	err := db.QueryRow(`SELECT UserID, FirstName, LastName, Email, Created FROM users WHERE UserID = ?`, usrID).Scan(&usr.ID, &usr.FirstName, &usr.LastName, &usr.Email, &usr.Created)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, errors.New("error occured getting user from db " + err.Error())
	}

	return &usr, nil
}

func LoginUser(usr customTypes.LoginUserRequest) (string, error) {
	query := `SELECT UserID, Password FROM users where email = ?`
	return LoginHelper(usr.Email, usr.Password, query)
}

func LoginAdmin(adm customTypes.LoginAdminRequest) (string, error) {
	query := `SELECT AdminID, Password FROM admins where email = ?`
	return LoginHelper(adm.Email, adm.Password, query)
}

func LoginHelper(email, password, query string) (string, error) {
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

func GetAdminByID(admID string) (*customTypes.Admin, error) {

	var adm customTypes.Admin

	err := db.QueryRow(`SELECT AdminID, Email, UserName, Created FROM admins WHERE AdminID = ?`, admID).Scan(&adm.ID, &adm.Email, &adm.UserName, &adm.Created)

	if err == sql.ErrNoRows {
		return nil, errors.New("admin not found")
	}

	if err != nil {
		return nil, errors.New("error occured getting admin from db " + err.Error())
	}

	return &adm, nil
}

func DeletePerson(person customTypes.Person, id string) error {

	var result sql.Result
	var err error

	switch person {
	case customTypes.USER:
		result, err = db.Exec(`DELETE FROM users WHERE UserID = ?`, id)
	case customTypes.ADMIN:
		result, err = db.Exec(`DELETE FROM admins WHERE AdminID = ?`, id)
	default:
		return errors.New("invalid person type")
	}

	if err != nil {
		return errors.New("error while deleting db " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return errors.New("error while checking affected rows: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return err
}

func SearchPersons(person customTypes.Person, usrRequest *customTypes.SearchUserRequest, _ *customTypes.SearchAdminRequest) (*[]customTypes.User, *[]customTypes.Admin, error) {
	var userList []customTypes.User

	var rows *sql.Rows
	var err error

	if person == customTypes.USER {
		rows, err = db.Query(`SELECT UserID, FirstName, LastName, Email, Created FROM users WHERE UserId = ? OR LOWER(FirstName) LIKE ? OR LOWER(LastName) LIKE ? OR LOWER(Email) LIKE ?`, usrRequest.ID, "%"+strings.ToLower(usrRequest.FirstName)+"%", "%"+strings.ToLower(usrRequest.LastName)+"%", "%"+strings.ToLower(usrRequest.Email)+"%")
	} else {
		return nil, nil, errors.New("invalid person type")
	}

	if err != nil {
		return nil, nil, errors.New("unable to perform query " + err.Error())
	}

	defer rows.Close()

	if person == customTypes.USER {
		for rows.Next() {
			var current customTypes.User

			err := rows.Scan(&current.ID, &current.FirstName, &current.LastName, &current.Email, &current.Created)

			if err != nil {
				return nil, nil, errors.New("error while appending users " + err.Error())
			}

			userList = append(userList, current)
		}

		return &userList, nil, nil
	}

	return nil, nil, errors.New("unable to perform query")
}

func AddAdmin(adm *customTypes.AddAdminRequest) (*customTypes.Admin, error) {
	var mail string

	err := db.QueryRow(`SELECT Email FROM admins where Email = ?`, adm.Email).Scan(&mail)

	if err == nil {
		return nil, errors.New("admin already exists")
	}

	if err != sql.ErrNoRows {
		return nil, errors.New("couldn't execute admin search in database: " + err.Error())
	}

	// create new admin
	var newAdmin customTypes.Admin
	var IDerr error
	newAdmin.ID, IDerr = uuid.NewUUID()

	if IDerr != nil {
		return nil, errors.New("couldn't generate UUID: " + IDerr.Error())
	}

	newAdmin.Created = int(time.Now().Unix())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adm.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("couldn't hash password: " + err.Error())
	}

	newAdmin.Password = string(hashedPassword)

	newAdmin.Email = adm.Email
	newAdmin.UserName = adm.UserName

	var rows *sql.Rows
	rows, err = db.Query(`INSERT INTO admins (AdminID, Email, Username, Password, Created) VALUES (?, ?, ?, ?, ?)`, newAdmin.ID, newAdmin.Email, newAdmin.UserName, newAdmin.Password, newAdmin.Created)

	if err != nil {
		return nil, errors.New("couldn't execute admin creation on db: " + err.Error())
	}

	defer rows.Close()

	fmt.Println("Server: New admin created: ID: ", newAdmin.ID)

	return &newAdmin, err
}
