package main

import (
	"encoding/json"
	"net/http"
)

// interface for api functions
type apiFunction func(http.ResponseWriter, *http.Request) error

func (server *Server) getBier(writer http.ResponseWriter, request *http.Request) error {
	return json.NewEncoder(writer).Encode("Bier")
}
