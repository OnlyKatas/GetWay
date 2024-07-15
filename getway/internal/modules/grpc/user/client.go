package user

import (
	"context"
	"fmt"
	"getway/internal/models"

	pbU "getway/pkg/proto/user/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

type Client struct {
	pbU.UserClient
}

func NewUserClient(ctx context.Context, addr string) (*Client, error) {
	const op = "pay.NewUserClient"

	clientConn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при установке соединения с сервером, %w", op, err)
	}

	return &Client{UserClient: pbU.NewUserClient(clientConn)}, nil
}

func (c Client) Get(id string) (*models.User, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", id, err)
	}
	req := &pbU.GetRequest{
		UserId: int64(intId),
	}
	resp, err := c.UserClient.Get(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", id, err)
	}

	user := models.User{
		ID:        int(resp.GetId()),
		FirstName: resp.GetFirstName(),
		LastName:  resp.GetLastName(),
		Username:  resp.GetUserName(),
		Email:     resp.GetEmail(),
		UUID:      resp.GetUuid(),
	}

	return &user, nil
}
