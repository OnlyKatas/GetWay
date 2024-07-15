package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"user/internal/models"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (u *UserStorage) AddUser(user *models.User) (int64, error) {
	query := `INSERT INTO users (uuid, first_name, last_name, username, email, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := u.db.QueryRow(query, user.UUID, user.FirstName, user.LastName, user.Username, user.Email, user.Password)
	var userId int
	if err := row.Scan(&userId); err != nil {
		return 0, err
	}
	return int64(userId), nil
}

func (u *UserStorage) GetUser(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, uuid, first_name, last_name, username FROM users WHERE id = $1`
	if err := u.db.Get(&user, query, id); err != nil {
		return nil, fmt.Errorf("user not found with id %d: %w", id, err)
	}

	return &user, nil
}

func (u *UserStorage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT password, id, uuid FROM users WHERE email = $1`
	if err := u.db.Get(&user, query, email); err != nil {
		return nil, fmt.Errorf("user not found with email %s: %w", email, err)
	}
	return &user, nil
}
