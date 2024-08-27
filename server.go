package main

import (
	"fmt"
	"net/http"
	"time"

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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, xJwtToken, ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (server *Server) run() {
	router := mux.NewRouter()

	// use CORS middleware to allow cross domain requests, fix later whith nginx oder some other shit
	router.Use(corsMiddleware)

	/*
		register routes

		note: api-function gets passed through error handling function
		to handle errors locally

	*/
	router.HandleFunc("/bier", handleError(server.handleGetBier)).Methods("GET", "OPTIONS")
	router.HandleFunc("/register", handleError(server.handleRegisterUser)).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", handleError(server.handleLoginUser)).Methods("POST", "OPTIONS")

	/*
		guarded api routes
	*/

	router.HandleFunc("/user/{ID}", JWTAuth(handleError(server.handleGetUserByID))).Methods("GET", "OPTIONS")

	/*
		admin routes for dashboard
	*/

	router.HandleFunc("/admin/login", handleError(server.handleLoginAdmin)).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin/validateJWT", handleError(server.handleValidateAdminJWT)).Methods("POST", "OPTIONS")

	/*
		guarded admin api routes
	*/

	router.HandleFunc("/admin/{ID}", JWTAuth(handleError(server.handleGetAdminByID))).Methods("GET", "OPTIONS")
	router.HandleFunc("/admins", JWTAuth(handleError(server.hanldeGetMultibleAdmins))).Methods("GET", "OPTIONS")
	router.HandleFunc("/admins/edit/{ID}", JWTAuth(handleError(server.handleEditAdmin))).Methods("POST", "OPTIONS")

	fmt.Println("Server: Running and Listening on port: ", server.adress)

	serverhandler := &http.Server{
		Addr:         server.adress,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err := serverhandler.ListenAndServe()

	if err != nil {
		fmt.Printf("Server: Error running server: %v\n", err)
	}
}
