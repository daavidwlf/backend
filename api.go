package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
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

	return json.NewDecoder(request.Body).Decode(content)
}

// function to writer error in a consistent format
func writeError(writer http.ResponseWriter, statusCode int, errmsg error) {
	err := writeJSON(writer, statusCode, map[string]string{"message": errmsg.Error()})
	if err != nil {
		fmt.Println("Server: Error ocurred: ", err.Error())
	}
}

// middleware
func JWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		tokenString := request.Header.Get("X-JWT-Token")

		token, err := validateJWT(tokenString)

		if err != nil {
			err := writeJSON(writer, http.StatusForbidden, map[string]string{"message": "permission denied"})
			if err != nil {
				fmt.Println("Server: Error ocurred: ", err.Error())
			}
			return
		}

		if !token.Valid {
			err := writeJSON(writer, http.StatusForbidden, map[string]string{"message": "permission denied"})
			if err != nil {
				fmt.Println("Server: Error ocurred: ", err.Error())
			}
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		reqID := request.Header.Get("ID")

		if reqID != claims["ID"] {
			err := writeJSON(writer, http.StatusForbidden, map[string]string{"message": "invalid token"})
			if err != nil {
				fmt.Println("Server: Error ocurred: ", err.Error())
			}
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

func createJWT(usrID string) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"ID":        usrID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

/*

	functions to handle routes

*/

func (server *Server) handleGetBier(writer http.ResponseWriter, _ *http.Request) error {
	return writeJSON(writer, http.StatusOK, "Bier")
}

func (server *Server) handleRegisterUser(writer http.ResponseWriter, request *http.Request) error {
	var userStruct registerUserRequest
	err := parseJSON(request, &userStruct)

	if err != nil {
		return err
	}

	err = registerUser(userStruct)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully created user"})
}

func (server *Server) handleLoginUser(writer http.ResponseWriter, request *http.Request) error {
	var usr loginUserRequest
	err := parseJSON(request, &usr)

	if err != nil {
		return err
	}

	var usrID string
	usrID, err = loginUser(usr)

	if err != nil {
		return err
	}

	// create jwt token when user logs in
	tokenString, err := createJWT(usrID)

	if err != nil {
		return errors.New("error while creating jwt token uuid: " + err.Error())
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully Logged in", "X-JWT-Token": tokenString})
}

func (server *Server) handleGetUserByID(writer http.ResponseWriter, request *http.Request) error {
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

func (server *Server) handleLoginAdmin(writer http.ResponseWriter, request *http.Request) error {
	var adm loginAdminRequest
	err := parseJSON(request, &adm)

	if err != nil {
		return err
	}

	var admID string
	admID, err = loginAdmin(adm)

	if err != nil {
		return err
	}

	// create jwt token when admin logs in
	tokenString, err := createJWT(admID)

	if err != nil {
		return errors.New("error while creating jwt token uuid: " + err.Error())
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully Logged in", "X-JWT-Token": tokenString, "adminID": admID})
}

func (server *Server) handleValidateAdminJWT(writer http.ResponseWriter, request *http.Request) error {
	var jwtRequest validateJWTRequest
	err := parseJSON(request, &jwtRequest)

	if err != nil {
		err := writeJSON(writer, http.StatusForbidden, map[string]string{"message": "Ainvalid token" + err.Error()})
		if err != nil {
			fmt.Println("Server: Error ocurred: ", err.Error())
		}
		return nil
	}

	var token *jwt.Token
	token, err = validateJWT(jwtRequest.Token)

	if err != nil {
		writeJSON(writer, http.StatusForbidden, map[string]string{"message": "Binvalid token"})
		return nil
	}

	if !token.Valid {
		writeJSON(writer, http.StatusForbidden, map[string]string{"message": "Cinvalid token"})
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)

	if jwtRequest.ID != claims["ID"] {
		writeJSON(writer, http.StatusForbidden, map[string]string{"message": "Dinvalid token"})
		return nil
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"message": "valid token"})
}

func (server *Server) handleGetAdminByID(writer http.ResponseWriter, request *http.Request) error {
	reqID := mux.Vars(request)["ID"]

	if reqID == "" {
		return errors.New("asinvalid ID")
	}

	adm, err := getAdminByID(reqID)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, adm)
}

func (server *Server) hanldeGetMultibleAdmins(writer http.ResponseWriter, request *http.Request) error {
	quantityParam := request.URL.Query().Get("quantity")

	var quantity int
	var err error
	if quantityParam != "" {
		quantity, err = strconv.Atoi(quantityParam)

		if err != nil {
			return errors.New("unable to parse quantity")
		}
	} else {
		quantity = 10
	}
	var adminList *[]admin
	adminList, err = getMultibleAdmins(quantity)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, adminList)
}

func (server *Server) handleEditAdmin(writer http.ResponseWriter, request *http.Request) error {

	adminID := mux.Vars(request)["ID"]

	if adminID == "" {
		return errors.New("id invalid")
	}

	var editAdm editAdminRequest

	err := parseJSON(request, &editAdm)

	if err != nil {
		return errors.New("unable to parse json" + err.Error())
	}

	var adm *editAdminRequest

	adm, err = editAdmin(adminID, &editAdm)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, adm)
}
