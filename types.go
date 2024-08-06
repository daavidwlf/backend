package main

import "github.com/google/uuid"

/*
*	variables must be in CAPTIAL letters otherwise they won't be exported when marashalling json!!!
 */

type festival struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	//to be completed
}

func newFestival(name string) festival {
	return festival{
		Id:   uuid.New(),
		Name: name,
	}
}

type user struct {
	Id              uuid.UUID
	Surname         string
	Name            string
	Email           string
	Pw_Endoded      string
	Created         int
	Saved_Festivals []festival
	//to be completed
}

func newUser(surname string, name string, email string, pw_endoced string, created int) *user {
	return &user{
		Id:         uuid.New(),
		Surname:    surname,
		Name:       name,
		Email:      email,
		Pw_Endoded: pw_endoced,
		Created:    created,
	}
}
