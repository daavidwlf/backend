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
