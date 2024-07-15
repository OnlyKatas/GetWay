package models

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
