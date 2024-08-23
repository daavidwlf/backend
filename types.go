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
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Created   int       `json:"created"`
}
