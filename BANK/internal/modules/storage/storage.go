package storage

import (
	"bank/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"time"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) SendDepositToDB(request *models.Request) error {
	query := `INSERT INTO bank (user_uuid, bank_uuid, amount, status) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, request.UserUUID, request.BankUUID, request.Amount, 0)
	if err != nil {
		return fmt.Errorf("error sending bank session to db %v", err)
	}
	return nil
}

func (s Storage) CanceledOP(bankUUID string) {
	time.Sleep(5 * time.Second)
	status := rand.Intn(2)
	if status == 0 {
		status = 1
	}
	query := `UPDATE bank SET status = $1 WHERE bank_uuid = $2`
	_, err := s.db.Exec(query, status, bankUUID)
	if err != nil {
		log.Printf("cancededOP.%v", err)
	}
}

func (s Storage) ReqForStatus(bankUUID, userUUID string) (int, error) {
	var status int
	query := `SELECT status FROM bank WHERE bank_uuid = $1 AND user_uuid = $2`
	if err := s.db.Get(status, query, bankUUID, userUUID); err != nil {
		return -1, fmt.Errorf("error getting bank status %v", err)
	}
	return status, nil
}
