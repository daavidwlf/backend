package main

import (
	"encoding/json"
	"errors"
	"net/http"

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
	var userStruct user
	err := parseJSON(request, &userStruct)

	if err != nil {
		return err
	}

	err = loginUser(userStruct)

	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, map[string]string{"Message": "Sucessfully Logged in"})
}

func (server *Server) getUserID(writer http.ResponseWriter, request *http.Request) error {
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
