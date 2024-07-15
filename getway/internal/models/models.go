package models

import "github.com/go-chi/jwtauth"

type User struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Credentials struct {
	ID     int    `json:"id"`
	UUID   string `json:"uuid"`
	Amount int    `json:"amount"`
}

type Token struct {
	Token string `json:"token"`
}

var TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
