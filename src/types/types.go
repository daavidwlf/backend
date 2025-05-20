package customTypes

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ApiFunction func(http.ResponseWriter, *http.Request) error

type Server struct {
	Adress string
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"userId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Created   int       `json:"created"`
}

type LoginAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ValidateJWTRequest struct {
	Token string `json:"xJwtToken"`
	ID    string `json:"adminId"`
}

type Admin struct {
	ID       uuid.UUID `json:"adminId"`
	UserName string    `json:"userName"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Created  int       `json:"created"`
}

type EditAdminRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type AddAdminRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DockerContainer struct {
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

type EditUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type SearchUserRequest struct {
	// ID this is a string so i won't throw an parse error when not searching with valid id
	ID        string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type SearchAdminRequest struct {
	// ID this is a string so i won't throw an parse error when not searching with valid id
	ID       string `json:"userId"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type Person int

const (
	USER Person = iota
	ADMIN
)

type LoginAttemptInfo struct {
	AttemptCount int
	LastAttempt  time.Time
	BlockedUntil time.Time
	IpAttempts   map[string]int
}
