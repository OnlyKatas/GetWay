package service

import (
	"fmt"
	"payS/internal/models"
)

type Service struct {
	Storage       Storageer
	SessionClinet SessionClinet
}

type Storageer interface {
	Hold(dto *models.DepositDTO) (*models.Session, error)
	Session(sessionID, userID, status int, bankID string) error
	Create(dto *models.DepositDTO) error
}

type SessionClinet interface {
	SendSession(s *models.Session) (string, error)
}

func NewService(storage Storageer, SessionClinet SessionClinet) *Service {
	return &Service{
		Storage:       storage,
		SessionClinet: SessionClinet,
	}
}

func (s *Service) Deposit(DTO *models.DepositDTO) (string, error) {

	session, err := s.Storage.Hold(DTO)
	if err != nil {
		return "", fmt.Errorf("не удалось создать сессию: %w", err)
	}
	res, err := s.SessionClinet.SendSession(session)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *Service) Create(DTO *models.DepositDTO) error {
	return s.Storage.Create(DTO)
}

func (s *Service) SetSessionStatus(sessionID, userID, status int, bankID string) error {
	return s.Storage.Session(sessionID, userID, status, bankID)
}
