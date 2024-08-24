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
			writeError(writer, http.StatusBadRequest, err)
		}
	}
}

func (server *Server) run() {
	router := mux.NewRouter()

	/*
		register routes

		note: api-function gets passed through error handling function
		to handle errors locally

	*/
	router.HandleFunc("/bier", handleError(server.handleGetBier)).Methods("GET")
	router.HandleFunc("/register", handleError(server.handleRegisterUser)).Methods("POST")
	router.HandleFunc("/login", handleError(server.handleLoginUser)).Methods("POST")

	/*
		guarded api routes
	*/

	router.HandleFunc("/user/{ID}", JWTAuth(handleError(server.handleGetUserByID))).Methods("GET")

	/*
		admin routes for dashboard
	*/

	router.HandleFunc("/admin/login", handleError(server.handleLoginAdmin)).Methods("POST")

	fmt.Println("Server: Running and Listening on port: ", server.adress)

	http.ListenAndServe(server.adress, router)
}
