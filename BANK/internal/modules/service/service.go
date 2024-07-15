package service

import (
	"bank/internal/models"
	"github.com/google/uuid"
)

type Service struct {
	storage PostgresAdapter
}

type PostgresAdapter interface {
	SendDepositToDB(request *models.Request) error
	CanceledOP(bankUUID string)
	ReqForStatus(bankUUID, userUUID string) (int, error)
}

func NewService(storage PostgresAdapter) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) ReqForPayProcessing(req *models.Request) (string, error) {
	bankSessionUUID := uuid.New().String()
	req.BankUUID = bankSessionUUID
	if err := s.storage.SendDepositToDB(req); err != nil {
		return "", err
	}
	go s.storage.CanceledOP(req.BankUUID)
	return bankSessionUUID, nil
}

func (s Service) ReqForPayStatus(req *models.Request) (int, error) {
	return s.storage.ReqForStatus(req.BankUUID, req.UserUUID)
}
