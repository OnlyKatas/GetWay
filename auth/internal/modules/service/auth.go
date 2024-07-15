package service

import (
	global "auth/internal/models"
	"context"
	"errors"

	"auth/internal/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB UserCreatorAndFinderService
}

type UserCreatorAndFinderService interface {
	Add(user *global.User) (int, string, error) /*Добавляет пользователя в БД*/
	Login(email string) (*global.User, error)   /*Ищет пользователя в БД, чтобы провести аутентификацию*/
}

func NewAuthService(u UserCreatorAndFinderService) *AuthService {
	return &AuthService{
		DB: u,
	}
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (token string, err error) {
	user, err := a.DB.Login(email)
	if err != nil {
		return "", errors.New("Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("Invalid username or password")
	}

	// Генерируем токен с ID пользователя
	_, tokenString, err := models.TokenAuth.Encode(jwt.MapClaims{
		"email": email,
		"id":    user.ID, // Добавляем ID пользователя в claims
		"uuid":  user.UUID,
	})
	if err != nil {
		return "", errors.New("Failed to generate token")
	}
	return tokenString, nil
}

func (a *AuthService) Register(ctx context.Context, user *global.User) (userID int64, uuid string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", err
	}
	user.Password = string(hashedPassword)

	id, uuid, err := a.DB.Add(user)
	if err != nil {
		return 0, "", err
	}
	return int64(id), uuid, nil
}
