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
	ID        uuid.UUID `json:"userID"`
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
	Token string `json:"x-jwt-token"`
	ID    string `json:"adminID"`
}

type admin struct {
	ID       uuid.UUID `json:"adminID"`
	UserName string    `json:"userName"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Created  int       `json:"created"`
}
