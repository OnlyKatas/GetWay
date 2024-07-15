package auth

import (
	"context"
	"encoding/json"
	"getway/internal/models"

	"log"
	"strconv"

	"net/http"
)

type AuthController struct {
	AuthServicer
	UserPaymentCreator
}

type AuthServicer interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user *models.User) (int64, string, error)
}

type UserPaymentCreator interface {
	Create(id int, uuid string) error
}

func NewAuthConroller(au AuthServicer, pc UserPaymentCreator) *AuthController {
	return &AuthController{
		AuthServicer:       au,
		UserPaymentCreator: pc,
	}
}

func (d *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	token, err := d.AuthServicer.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, "Failed to encode token", http.StatusInternalServerError)
	}
}
func (d *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id, uuid, err := d.Register(r.Context(), &input)
	if err != nil {
		http.Error(w, "Failed to register pay", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := d.Create(int(id), uuid); err != nil {
		http.Error(w, "Failed to create pay account", http.StatusInternalServerError)
		log.Println(err)
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"id": strconv.Itoa(int(id))}); err != nil {
		http.Error(w, "Failed to encode id", http.StatusInternalServerError)
	}
}
