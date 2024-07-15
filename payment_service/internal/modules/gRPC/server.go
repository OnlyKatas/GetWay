package controller

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"payS/internal/models"

	protobuff "payS/pkg/proto/pay/gen/go"
)

type Payer interface {
	Deposit(DTO *models.DepositDTO) (string, error)
	Create(DTO *models.DepositDTO) error
	SetSessionStatus(sessionID, userID, status int, bankID string) error
}

type Server struct {
	protobuff.PayServer
	service Payer
}

func Register(gRPC *grpc.Server, service Payer) {
	protobuff.RegisterPayServer(gRPC, &Server{service: service})
}

func NewServer(service Payer) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Deposit(ctx context.Context, req *protobuff.DepositRequest) (resp *protobuff.DepositResponse, err error) {
	DTO := models.DepositDTO{
		UserID:   int(req.GetUserId()),
		UserUUID: req.GetUuid(),
		Amount:   int(req.GetAmount()),
	}
	answer, err := s.service.Deposit(&DTO)
	if err != nil {
		return nil, fmt.Errorf("server.depost failed to deposit: %w", err)
	}
	resp = &protobuff.DepositResponse{
		Req: answer,
	}
	return resp, nil
}

func (s *Server) Create(ctx context.Context, req *protobuff.CreateRequest) (*emptypb.Empty, error) {
	DTO := models.DepositDTO{
		UserID:   int(req.GetUserId()),
		UserUUID: req.GetUuid(),
	}
	if err := s.service.Create(&DTO); err != nil {
		return nil, fmt.Errorf("server.create: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) SetSessionStatus(ctx context.Context, req *protobuff.SetRequest) (*emptypb.Empty, error) {
	if err := s.service.SetSessionStatus(int(req.GetSessionId()), int(req.GetUserId()), int(req.GetSessionStatus()), req.GetBankId()); err != nil {
		return nil, fmt.Errorf("server.setSessionStatus: %w", err)
	}
	return &emptypb.Empty{}, nil
}
