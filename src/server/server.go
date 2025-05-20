package server

import (
	"backend/src/api"
	customTypes "backend/src/types"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CreateServer(adress string) *customTypes.Server {
	return &customTypes.Server{
		Adress: adress,
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

func Run(s *customTypes.Server) {
	router := mux.NewRouter()

	// use CORS middleware to allow cross domain requests, fix later whith nginx oder some other shit
	router.Use(corsMiddleware)

	/*
		register routes

		note: api-function gets passed through error handling function
		to handle errors locally

	*/
	router.HandleFunc("/bier", api.HandleError(api.HandleGetBier)).Methods("GET", "OPTIONS")
	router.HandleFunc("/register", api.HandleError(api.HandleRegisterUser)).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", api.HandleError(api.HandleLoginUser)).Methods("POST", "OPTIONS")

	/*
		guarded api routes
	*/

	router.HandleFunc("/user/{ID}", api.JWTAuth(api.HandleError(api.HandleGetUserByID))).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", api.JWTAuth(api.HandleError(api.HandleGetMultibleUsers))).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/search", api.JWTAuth(api.HandleError(api.HandleSearchUsers))).Methods("POST", "OPTIONS")

	router.HandleFunc("/user/edit/{ID}", api.JWTAuth(api.HandleError(api.HandleEditUser))).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/delete/{ID}", api.JWTAuth(api.HandleError(api.HandleDeleteUser))).Methods("POST", "OPTIONS")

	/*
		admin routes for dashboard
	*/

	router.HandleFunc("/admin/login", api.HandleError(api.HandleLoginAdmin)).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin/validateJWT", api.HandleError(api.HandleValidateAdminJWT)).Methods("POST", "OPTIONS")

	/*
		guarded admin api routes
	*/

	router.HandleFunc("/admin/{ID}", api.JWTAuth(api.HandleError(api.HandleGetAdminByID))).Methods("GET", "OPTIONS")
	router.HandleFunc("/admins", api.JWTAuth(api.HandleError(api.HandleGetMultibleAdmins))).Methods("GET", "OPTIONS")

	router.HandleFunc("/admin/edit/{ID}", api.JWTAuth(api.HandleError(api.HandleEditAdmin))).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin/delete/{ID}", api.JWTAuth(api.HandleError(api.HandleDeleteAdmin))).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin/add", api.JWTAuth(api.HandleError(api.HandleAddAdmin))).Methods("POST", "OPTIONS")

	router.HandleFunc("/docker/containers", api.JWTAuth(api.HandleError(api.HandleGetDockerContainers))).Methods("GET", "OPTIONS")

	fmt.Println("Server: Running and Listening on port: ", s.Adress)

	serverhandler := &http.Server{
		Addr:         s.Adress,
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
