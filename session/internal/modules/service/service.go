package service

import (
	"fmt"
	"session/internal/models"
)

type Service struct {
	BankAPI BankAPI
}

type BankAPI interface {
	Deposit(userUUID string, amount int) error
}

func NewService(bank BankAPI) *Service {
	return &Service{
		BankAPI: bank,
	}
}

func (s *Service) SendSessionToBank(ses *models.Session) (string, error) {
	if err := s.BankAPI.Deposit(ses.UserUUID, ses.Amount); err != nil {
		return "", fmt.Errorf("failed to deposit session to bank api: %w", err)
	}
	return fmt.Sprint("Операция успешно отправлена банк, ожидайте"), nil
}
