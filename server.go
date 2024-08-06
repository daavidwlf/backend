package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	adress string
}

func createServer(adress string) *Server {
	return &Server{
		adress: adress,
	}
}

func handleError(function apiFunction) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := function(writer, request)
		if err != nil {
			fmt.Println("Server: Error ocurred: ", err.Error())
		}
	}
}

func (server *Server) run() {
	router := mux.NewRouter()

	/*
	*	register routes
	*
	*	note: apiFunction gets passed through error handling function
	*	to handle errors locally
	*
	 */
	router.HandleFunc("/bier", handleError(server.getBier))

	fmt.Println("Server: Running and Listening on port: ", server.adress)

	http.ListenAndServe(server.adress, router)
}
