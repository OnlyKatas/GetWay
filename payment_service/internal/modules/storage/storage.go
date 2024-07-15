package storage

import (
	"fmt"
	"log"
	"payS/internal/models"

	"github.com/jmoiron/sqlx"
)

// Storage struct to manage database operations
type Storage struct {
	db *sqlx.DB
}

// NewStorage initializes a new Storage instance
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Hold adds a deposit to user_balance and creates a payment session
func (s *Storage) Hold(dto *models.DepositDTO) (*models.Session, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("не могу начать транзакцию: %w", err)
	}

	// Update user hold balance
	updateHoldBalanceQuery := `UPDATE user_balance SET hold_balance = hold_balance + $1 WHERE user_id = $2 AND user_uuid = $3`
	_, err = tx.Exec(updateHoldBalanceQuery, dto.Amount, dto.UserID, dto.UserUUID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("не могу обновить hold баланс пользователя: %w", err)
	}

	// Create payment session
	createSessionQuery := `INSERT INTO payment_session (user_id, user_uuid, amount, status) VALUES ($1, $2, $3, $4) RETURNING id`
	row := tx.QueryRow(createSessionQuery, dto.UserID, dto.UserUUID, dto.Amount, 1)
	var sessionID int
	err = row.Scan(&sessionID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("не могу создать платежную сессию: %w", err)
	}
	session := models.Session{
		ID:       sessionID,
		UserID:   dto.UserID,
		UserUUID: dto.UserUUID,
		Amount:   dto.Amount,
	}

	tx.Commit()
	log.Println("Транзакция завершилась успешно\n", dto)
	return &session, nil
}

func (s *Storage) Session(sessionID, userID, status int, bankID string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("не могу начать транзакцию: %w", err)
	}
	log.Println("транзакция началась")

	var session models.Session
	query := `SELECT id, uuid, user_id, user_uuid, amount, created_at, updated_at, status FROM payment_session WHERE id = $1 AND user_id = $2`
	err = tx.Get(&session, query, sessionID, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("не могу получить платежную сессию: %w", err)
	}

	if status == 0 {
		// Update real balance and hold balance
		updateBalanceQuery := `UPDATE user_balance SET balance = balance + $1, hold_balance = hold_balance - $1 WHERE user_id = $2 AND user_uuid = $3`
		_, err := tx.Exec(updateBalanceQuery, session.Amount, session.UserID, session.UserUUID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не могу обновить баланс пользователя: %w", err)
		}
	} else if status == 3 {
		// Reduce hold balance only
		updateHoldBalanceQuery := `UPDATE user_balance SET hold_balance = hold_balance - $1 WHERE user_id = $2 AND user_uuid = $3`
		_, err := tx.Exec(updateHoldBalanceQuery, session.Amount, session.UserID, session.UserUUID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не могу обновить hold баланс пользователя: %w", err)
		}
		log.Println("Банк отклонил платеж для сессии:", session.ID)
	}

	updateBankIDQuery := `UPDATE payment_session SET bank_id = $1 WHERE id = $2`
	_, err = tx.Exec(updateBankIDQuery, bankID, sessionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("не могу обновить bankID в платежной сессии: %w", err)
	}

	tx.Commit()
	log.Println("транзакция прошла успешно")
	return nil
}

// Create adds a new user to user_balance table
func (s *Storage) Create(dto *models.DepositDTO) error {
	query := `INSERT INTO user_balance (user_id, user_uuid, balance, hold_balance) VALUES ($1, $2, 0, 0)`
	_, err := s.db.Exec(query, dto.UserID, dto.UserUUID)
	if err != nil {
		return fmt.Errorf("не могу добавить пользователя: %w", err)
	}
	log.Println("Пользователь успешно создан", dto)
	return nil
}
