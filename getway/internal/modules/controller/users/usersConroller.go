package users

import (
	"encoding/json"
	"getway/internal/models"

	"log"
	"net/http"
)

type UserController struct {
	UserServicer
}

func NewUsersController(u UserServicer) *UserController {
	return &UserController{u}
}

type UserServicer interface {
	Get(id string) (*models.User, error)
}

func (u *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	id := vars.Get("id")
	if id == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	answer, err := u.UserServicer.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}
