package service

import (
	"context"
	"user/internal/models"

	"github.com/google/uuid"
)

type UserService struct {
	US UserStorageer
}

type UserStorageer interface {
	AddUser(user *models.User) (int64, error)
	GetUser(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

func NewUserService(us UserStorageer) *UserService {
	return &UserService{
		US: us,
	}
}

func (u *UserService) Get(ctx context.Context, id int64) (*models.User, error) {
	return u.US.GetUser(int(id))
}

func (u *UserService) Add(ctx context.Context, user *models.User) (int64, string, error) {
	user.UUID = uuid.New().String()
	id, err := u.US.AddUser(user)
	if err != nil {
		return 0, "", err
	}
	return id, user.UUID, nil
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.US.GetUserByEmail(email)
}
