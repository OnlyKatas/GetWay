package pay

import (
	"encoding/json"
	"fmt"
	"getway/internal/models"
	"github.com/go-chi/jwtauth"
	"log"
	"net/http"
)

type Controller struct {
	PayServicer
}

type PayServicer interface {
	Deposit(id int, uuid string, amount int) (string, error)
	Create(id int, uuid string) error
}

func NewController(p PayServicer) *Controller {
	return &Controller{
		PayServicer: p,
	}
}

func (c *Controller) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	var (
		input  models.Credentials
		err    error
		answer string
	)
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	input.ID, input.UUID, err = getIDAndUUID(claims)
	if err != nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	answer, err = c.PayServicer.Deposit(input.ID, input.UUID, input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error calling PayServicer.Deposit:", err)
		return
	}

	if _, err = fmt.Fprint(w, answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error writing response:", err)
	}
}

func getIDAndUUID(claims map[string]interface{}) (int, string, error) {
	// Разбор и проверка токена
	idFloat, ok := claims["id"].(float64) // JWT может содержать числа в формате float64
	if !ok {
		return 0, "", fmt.Errorf("не могу достать id: %v", claims["id"])
	}
	id := int(idFloat) // Преобразование float64 в int

	uuid, ok := claims["uuid"].(string)
	if !ok {
		return 0, "", fmt.Errorf("не могу достать uuid: %v", claims["uuid"])
	}
	return id, uuid, nil
}
