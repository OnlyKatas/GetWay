package auth

import (
	"context"
	"fmt"
	"getway/internal/models"

	pb "getway/pkg/proto/auth/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	authClient pb.AuthClient
}

// NewAuthService создаёт новый клиент Client
func NewAuthClient(ctx context.Context, addr string) (*Client, error) {
	const op = "auth.NewAuthService"

	clientConn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{authClient: pb.NewAuthClient(clientConn)}, nil

}

// Login выполняет gRPC вызов Login
func (s *Client) Login(ctx context.Context, email, password string) (string, error) {
	req := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}

	// Установим таймаут для контекста
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	res, err := s.authClient.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}

	return res.GetToken(), nil
}

// Register выполняет gRPC вызов Register
func (s *Client) Register(ctx context.Context, user *models.User) (int64, string, error) {
	req := &pb.RegisterRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
	}

	// Установим таймаут для контекста
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	res, err := s.authClient.Register(ctx, req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to register: %w", err)
	}

	return res.GetUserId(), res.GetUuid(), nil
}
