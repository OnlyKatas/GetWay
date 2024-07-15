package payGRPC

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"session/internal/models"
	pb "session/pkg/proto/pay/gen/go"
)

type Client struct {
	payClient pb.PayClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	const op = "payGRPC.NewClient"
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Client{payClient: pb.NewPayClient(conn)}, nil
}

func (c *Client) GetSessionStatus(ctx context.Context) ([]models.Session, error) {
	res, err := c.payClient.GetSessionStatus(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("grpc get session status: %w", err)
	}
	sessions := make([]models.Session, 0, len(res.Sessions))
	for _, session := range res.GetSessions() {
		sessions = append(sessions, models.Session{
			ID:       int(session.GetId()),
			UUID:     session.GetUuid(),
			UserID:   int(session.GetUserId()),
			UserUUID: session.GetUserUuid(),
			Amount:   int(session.GetAmount()),
		})
	}
	return sessions, nil
}
