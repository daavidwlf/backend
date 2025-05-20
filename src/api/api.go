package api

import (
	"backend/src/db"
	customTypes "backend/src/types"
	"backend/src/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// function to write JSON
func WriteJSON(writer http.ResponseWriter, statusCode int, content any) error {
	/*
		This order of commands is mandatory!
	*/
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(content)
}

// function to parse JSONs
func ParseJSON(request *http.Request, content any) error {
	if request.Body == nil {
		return errors.New("body of request is nil")
	}

	return json.NewDecoder(request.Body).Decode(content)
}

// function to writer error in a consistent format
func WriteError(writer http.ResponseWriter, statusCode int, errmsg error) {
	err := WriteJSON(writer, statusCode, map[string]string{"message": errmsg.Error()})
	if err != nil {
		fmt.Println("Server: Error ocurred: ", err.Error())
	}
}

func HandleError(function customTypes.ApiFunction) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := function(writer, request)
		if err != nil {
			fmt.Println("Server: Error ocurred: ", err.Error())
			WriteError(writer, http.StatusBadRequest, err)
		}
	}
}

// middleware
func JWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		tokenString := request.Header.Get("xJwtToken")

		token, err := ValidateJWT(tokenString)

		if err != nil {
			err := WriteJSON(writer, http.StatusForbidden, map[string]string{"message": "permission denied"})
			if err != nil {
				fmt.Println("Server: Error ocurred: ", err.Error())
			}
			return
		}

		if !token.Valid {
			err := WriteJSON(writer, http.StatusForbidden, map[string]string{"message": "permission denied"})
			if err != nil {
				fmt.Println("Server: Error ocurred: ", err.Error())
			}
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		reqID := request.Header.Get("ID")

		if reqID != claims["ID"] {
			err := WriteJSON(writer, http.StatusForbidden, map[string]string{"message": "invalid token"})
			if err != nil {
				fmt.Println("Server: Error ocurred: ", err.Error())
			}
			return
		}

		handlerFunc(writer, request)
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func CreateJWT(usrID string) (string, error) {
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

func HandleGetBier(writer http.ResponseWriter, _ *http.Request) error {
	return WriteJSON(writer, http.StatusOK, "Bier")
}

func HandleRegisterUser(writer http.ResponseWriter, request *http.Request) error {
	var userStruct customTypes.RegisterUserRequest
	err := ParseJSON(request, &userStruct)

	if err != nil {
		return err
	}

	err = db.RegisterUser(userStruct)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully created user"})
}

func HandleLoginUser(writer http.ResponseWriter, request *http.Request) error {
	var usr customTypes.LoginUserRequest
	err := ParseJSON(request, &usr)

	if err != nil {
		return err
	}

	var usrID string
	usrID, err = db.LoginUser(usr)

	if err != nil {
		return err
	}

	// create jwt token when user logs in
	tokenString, err := CreateJWT(usrID)

	if err != nil {
		return errors.New("error while creating jwt token uuid: " + err.Error())
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully Logged in", "X-JWT-Token": tokenString})
}

func HandleEditUser(writer http.ResponseWriter, request *http.Request) error {
	userID := mux.Vars(request)["ID"]

	if userID == "" {
		return errors.New("id invalid")
	}

	var editUsr customTypes.EditUserRequest

	err := ParseJSON(request, &editUsr)

	if err != nil {
		return errors.New("unable to parse json" + err.Error())
	}

	var usrID string

	usrID, err = db.EditPerson(customTypes.USER, userID, &editUsr, nil)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully updated user  " + usrID})
}

func HandleGetUserByID(writer http.ResponseWriter, request *http.Request) error {
	reqID := mux.Vars(request)["ID"]

	if reqID == "" {
		return errors.New("invalid ID")
	}

	usr, err := db.GetUserByID(reqID)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, usr)
}

func HandleLoginAdmin(writer http.ResponseWriter, request *http.Request) error {

	var adm customTypes.LoginAdminRequest
	err := ParseJSON(request, &adm)

	if err != nil {
		return err
	}

	ip := request.RemoteAddr

	blocked := utils.TrackLoginAttempt(ip, adm.Email)

	if blocked {
		return errors.New("too many requests")
	}

	var admID string
	admID, err = db.LoginAdmin(adm)

	if err != nil {
		return err
	}

	// create jwt token when admin logs in
	tokenString, err := CreateJWT(admID)

	if err != nil {
		return errors.New("error while creating jwt token uuid: " + err.Error())
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully Logged in", "xJwtToken": tokenString, "adminId": admID})
}

func HandleValidateAdminJWT(writer http.ResponseWriter, request *http.Request) error {
	var jwtRequest customTypes.ValidateJWTRequest
	err := ParseJSON(request, &jwtRequest)

	if err != nil {
		err := WriteJSON(writer, http.StatusForbidden, map[string]string{"message": "Ainvalid token" + err.Error()})
		if err != nil {
			fmt.Println("Server: Error ocurred: ", err.Error())
		}
		return nil
	}

	var token *jwt.Token
	token, err = ValidateJWT(jwtRequest.Token)

	if err != nil {
		err := WriteJSON(writer, http.StatusForbidden, map[string]string{"message": "Binvalid token"})
		if err != nil {
			fmt.Println("Server: Error ocurred: ", err.Error())
		}
		return nil
	}

	if !token.Valid {
		err := WriteJSON(writer, http.StatusForbidden, map[string]string{"message": "Cinvalid token"})
		if err != nil {
			fmt.Println("Server: Error ocurred: ", err.Error())
		}
		return nil
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	if jwtRequest.ID != claims["ID"] {
		err := WriteJSON(writer, http.StatusForbidden, map[string]string{"message": "Dinvalid token"})
		if err != nil {
			fmt.Println("Server: Error ocurred: ", err.Error())
		}
		return nil
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "valid token"})
}

func HandleGetAdminByID(writer http.ResponseWriter, request *http.Request) error {
	reqID := mux.Vars(request)["ID"]

	if reqID == "" {
		return errors.New("asinvalid ID")
	}

	adm, err := db.GetAdminByID(reqID)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, adm)
}

func HandleGetMultibleUsers(writer http.ResponseWriter, request *http.Request) error {
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

	var userList *[]customTypes.User
	userList, _, err = db.GetMultiblePersons(customTypes.USER, quantity)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, userList)
}

func HandleDeleteUser(writer http.ResponseWriter, request *http.Request) error {
	userID := mux.Vars(request)["ID"]

	if userID == "" {
		return errors.New("id invalid")
	}

	err := db.DeletePerson(customTypes.USER, userID)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "user " + userID + " deleted"})
}

func HandleSearchUsers(writer http.ResponseWriter, request *http.Request) error {

	var userSearchRequest *customTypes.SearchUserRequest

	err := ParseJSON(request, &userSearchRequest)

	if err != nil {
		return errors.New("unable to parse json " + err.Error())
	}

	var userList *[]customTypes.User

	userList, _, err = db.SearchPersons(customTypes.USER, userSearchRequest, nil)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, userList)
}

func HandleGetMultibleAdmins(writer http.ResponseWriter, request *http.Request) error {
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
	var adminList *[]customTypes.Admin
	_, adminList, err = db.GetMultiblePersons(customTypes.ADMIN, quantity)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, adminList)
}

func HandleEditAdmin(writer http.ResponseWriter, request *http.Request) error {

	adminID := mux.Vars(request)["ID"]

	if adminID == "" {
		return errors.New("id invalid")
	}

	var editAdm customTypes.EditAdminRequest

	err := ParseJSON(request, &editAdm)

	if err != nil {
		return errors.New("unable to parse json" + err.Error())
	}

	var admID string

	admID, err = db.EditPerson(customTypes.ADMIN, adminID, nil, &editAdm)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "Sucessfully updated user  " + admID})
}

func HandleDeleteAdmin(writer http.ResponseWriter, request *http.Request) error {
	adminID := mux.Vars(request)["ID"]

	if adminID == "" {
		return errors.New("id invalid")
	}

	err := db.DeletePerson(customTypes.ADMIN, adminID)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "admin " + adminID + " deleted"})
}

func HandleAddAdmin(writer http.ResponseWriter, request *http.Request) error {

	var addAdm customTypes.AddAdminRequest

	err := ParseJSON(request, &addAdm)

	if err != nil {
		return errors.New("unable to parse json" + err.Error())
	}

	var newAdmin *customTypes.Admin

	newAdmin, err = db.AddAdmin(&addAdm)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, map[string]string{"message": "admin " + newAdmin.UserName + " successfullyy created"})
}

func HandleGetDockerContainers(writer http.ResponseWriter, _ *http.Request) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return errors.New("failed to create docker client: " + err.Error())
	}

	defer cli.Close()

	// set client version to docker deamon version
	cli.NegotiateAPIVersion(context.Background())

	runningContainers, err := cli.ContainerList(context.Background(), containerTypes.ListOptions{All: true})

	if err != nil {
		return errors.New("failed to list docker options: " + err.Error())
	}

	containers := make([]customTypes.DockerContainer, 0, len(runningContainers))

	for _, container := range runningContainers {
		var current customTypes.DockerContainer

		options := containerTypes.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     false,
			// only 50 newest lines get read out
			Tail: "50",
		}

		logs, err := cli.ContainerLogs(context.Background(), container.ID, options)

		if err != nil {
			return errors.New("unable to read conatiner logs: " + err.Error())
		}

		defer logs.Close()

		logsArray, err := utils.ParseLogs(logs)

		if err != nil {
			return errors.New("unable to parse logs: " + err.Error())
		}

		current.Logs = logsArray

		if len(container.Names) > 0 {
			current.Name = container.Names[0]
		}

		if len(container.Ports) > 0 {
			port := container.Ports[0]
			current.IP = port.IP
			current.PublicPort = port.PublicPort
			current.PrivatePort = port.PrivatePort
		}

		current.Created = container.Created

		current.State = container.State

		current.Status = container.Status

		current.Image = container.Image

		var volumes []string

		for _, mount := range container.Mounts {
			if mount.Type == "volume" {
				volumes = append(volumes, mount.Name)
			}
		}

		current.Volumes = volumes

		containers = append(containers, current)
	}

	return WriteJSON(writer, http.StatusOK, containers)
}
