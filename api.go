package main

import (
	"encoding/json"
	"net/http"
)

// interface for api functions
type apiFunction func(http.ResponseWriter, *http.Request) error

// function to write JSON for api functions
func writeJSON(writer http.ResponseWriter, content any) error {
	/*
	*	this order of commands is mandatory!!!
	 */
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	return json.NewEncoder(writer).Encode(content)
}

func (server *Server) getBier(writer http.ResponseWriter, request *http.Request) error {
	currentFestival := newFestival("Southside")
	return writeJSON(writer, currentFestival)
}
