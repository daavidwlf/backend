package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// type for api functions
type apiFunction func(http.ResponseWriter, *http.Request) error

// function to write JSON
func writeJSON(writer http.ResponseWriter, statusCode int, content any) error {
	/*
		This order of commands is mandatory!
	*/
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(content)
}

// function to parse JSONs
func parseJSON(request *http.Request, content any) error {
	if request.Body == nil {
		return errors.New("body of request is nil")
	}

	//fmt.Println(io.ReadAll(request.Body))

	return json.NewDecoder(request.Body).Decode(content)
}

// function to writer error in a consistent format
func writeError(writer http.ResponseWriter, statusCode int, err error) {
	writeJSON(writer, statusCode, map[string]string{"Error": err.Error()})
}

// middleware
func JWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		tokenString := request.Header.Get("X-JWT-Token")

		token, err := validateJWT(tokenString)

		if err != nil {
			writeJSON(writer, http.StatusForbidden, map[string]string{"Message": "permission denied"})
			return
		}

		if !token.Valid {
			writeJSON(writer, http.StatusForbidden, map[string]string{"Message": "permission denied"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		reqID := mux.Vars(request)["ID"]

		if reqID != claims["userID"] {
			writeJSON(writer, http.StatusForbidden, map[string]string{"Message": "invalid token"})
			return
		}

		handlerFunc(writer, request)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func createJWT(usr user) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"userID":    usr.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

/*

	functions to handle routes

*/

func (server *Server) getBier(writer http.ResponseWriter, request *http.Request) error {
	// currentFestival := newFestival("Southside")
	// return writeJSON(writer, http.StatusOK, currentFestival)
	return writeJSON(writer, http.StatusOK, "Bier")
}

func (server *Server) register(writer http.ResponseWriter, request *http.Request) error {
	//mandatory
	var userStruct user
	err := parseJSON(request, &userStruct)

	if err != nil {
		return err
	}

	err = registerUser(userStruct)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"Message": "Sucessfully created user"})
}

func (server *Server) login(writer http.ResponseWriter, request *http.Request) error {
	var usr user
	err := parseJSON(request, &usr)

	if err != nil {
		return err
	}

	var usrID string
	usrID, err = loginUser(usr)

	if err != nil {
		return err
	}

	usr.ID, err = uuid.Parse(usrID)

	if err != nil {
		return errors.New("error while parsing uuid: " + err.Error())
	}

	tokenString, err := createJWT(usr)

	if err != nil {
		return errors.New("error while creating jwt token uuid: " + err.Error())
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"Message": "Sucessfully Logged in", "X-JWT-Token": tokenString})
}

func (server *Server) getUserByID(writer http.ResponseWriter, request *http.Request) error {
	reqID := mux.Vars(request)["ID"]

	if reqID == "" {
		return errors.New("invalid ID")
	}

	usr, err := getUserByID(reqID)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, usr)
}
