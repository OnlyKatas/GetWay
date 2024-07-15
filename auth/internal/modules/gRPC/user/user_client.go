package user

import (
	global "auth/internal/models"
	PBUser "auth/pkg/proto/user/gen/go"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	PBUser.UserClient
}

func NewUserClient(addr string) (*UserClient, error) {
	const op = "user.NewUserClient"
	ctx := context.Background()
	clientConn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при установке соединения с сервером, %w", op, err)
	}

	return &UserClient{UserClient: PBUser.NewUserClient(clientConn)}, nil
}

func (c *UserClient) Login(email string) (*global.User, error) {

	req := &PBUser.LoginRequest{
		Email: email,
	}
	resp, err := c.UserClient.Login(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("Опреация логин не срабатывает: %w", err)
	}

	user := global.User{
		Password: resp.GetPassword(),
		ID:       int(resp.GetId()),
		UUID:     resp.GetUuid(),
	}
	return &user, nil
}

func (c *UserClient) Add(user *global.User) (int, string, error) {
	req := &PBUser.AddRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Uuid:      user.UUID,
	}
	resp, err := c.UserClient.Add(context.Background(), req)
	if err != nil {
		return 0, "", fmt.Errorf("Операция добавления новго пользователя не срабатывает: %w", err)
	}

	return int(resp.GetUserId()), resp.GetUuid(), nil
}
