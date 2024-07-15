package controller

import (
	"context"
	"fmt"
	"user/internal/models"

	"google.golang.org/grpc"
	protobuff "user/pkg/proto/user/gen/go"
)

type Userer interface {
	Get(ctx context.Context, id int64) (*models.User, error)
	Add(ctx context.Context, user *models.User) (int64, string, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type ServerAPI struct {
	protobuff.UserServer
	user Userer
}

func Register(gRPC *grpc.Server, user Userer) {
	protobuff.RegisterUserServer(gRPC, &ServerAPI{user: user})
}

func NewServerAPI(user Userer) *ServerAPI {
	return &ServerAPI{
		user: user,
	}
}

func (s *ServerAPI) Get(ctx context.Context, req *protobuff.GetRequest) (*protobuff.GetResponse, error) {
	user, err := s.user.Get(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	resp := &protobuff.GetResponse{
		Id:    int64(user.ID),
		Email: user.Email,
	}
	return resp, err
}

func (s *ServerAPI) Add(ctx context.Context, req *protobuff.AddRequest) (*protobuff.AddResponse, error) {
	user := models.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
	}
	id, uuid, err := s.user.Add(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("add user: %w", err)
	}
	return &protobuff.AddResponse{
		UserId: id,
		Uuid:   uuid,
	}, nil
}

func (s *ServerAPI) Login(ctx context.Context, req *protobuff.LoginRequest) (*protobuff.LoginResponse, error) {
	user, err := s.user.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	resp := &protobuff.LoginResponse{
		Email:    user.Email,
		Password: user.Password,
		Id:       int64(user.ID),
		Uuid:     user.UUID,
	}
	return resp, nil
}
