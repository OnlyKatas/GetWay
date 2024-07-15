package payGRPC

import (
	"context"
	"fmt"
	pb "getway/pkg/proto/pay/gen/go"
	"google.golang.org/grpc"
	"time"
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

func (c *Client) Deposit(id int, uuid string, amount int) (string, error) {
	req := &pb.DepositRequest{
		UserId: int64(id),
		Uuid:   uuid,
		Amount: int64(amount),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := c.payClient.Deposit(ctx, req)
	if err != nil {
		return "", fmt.Errorf("filed to call deposit method: %w", err)
	}

	return resp.GetReq(), nil
}

func (c *Client) Create(id int, uuid string) error {
	req := &pb.CreateRequest{
		UserId: int64(id),
		Uuid:   uuid,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := c.payClient.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("filed to call create method: %w", err)
	}
	return nil
}
