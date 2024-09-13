package main

import "github.com/google/uuid"

/*
	variables must start with CAPITIAL letters otherwise they won't be exported when marashalling json!
*/

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type user struct {
	ID        uuid.UUID `json:"userId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Created   int       `json:"created"`
}

type loginAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type validateJWTRequest struct {
	Token string `json:"xJwtToken"`
	ID    string `json:"adminId"`
}

type admin struct {
	ID       uuid.UUID `json:"adminId"`
	UserName string    `json:"userName"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Created  int       `json:"created"`
}

type editAdminRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type addAdminRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type dockerContainer struct {
	Name        string   `json:"name"`
	PublicPort  uint16   `json:"publicPort"`
	PrivatePort uint16   `json:"privatePort"`
	IP          string   `json:"ip"`
	Created     int64    `json:"created"`
	State       string   `json:"state"`
	Status      string   `json:"status"`
	Image       string   `json:"image"`
	Volumes     []string `json:"volume"`
	Logs        []string `json:"logs"`
}

type editUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type searchUserRequest struct {
	// ID this is a string so i won't throw an parse error when not searching with valid id
	ID        string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type searchAdminRequest struct {
	// ID this is a string so i won't throw an parse error when not searching with valid id
	ID       string `json:"userId"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type person int

const (
	USER person = iota
	ADMIN
)
