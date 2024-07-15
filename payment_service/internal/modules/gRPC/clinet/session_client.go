package clinet

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"payS/internal/models"
	pb "payS/pkg/proto/session/gen/go"
	"time"
)

type Client struct {
	payClient pb.SessionClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	const op = "payGRPC.NewClient"
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Client{payClient: pb.NewSessionClient(conn)}, nil
}

func (c *Client) SendSession(s *models.Session) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := c.payClient.SendSession(ctx, &pb.SendSessionRequest{
		Id:       int64(s.ID),
		UserId:   int64(s.UserID),
		UserUuid: s.UserUUID,
		Amount:   int64(s.Amount),
	})
	if err != nil {
		return "", fmt.Errorf("Ошибка на пути в банк или в банке: %w", err)
	}
	return res.GetReq(), nil
}
