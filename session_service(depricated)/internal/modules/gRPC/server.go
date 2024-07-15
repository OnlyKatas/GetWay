package controller

import (
	"google.golang.org/grpc"
	"session/internal/models"

	protobuff "session/pkg/proto/pay/gen/go"
)

type Uniplemented interface {
	Deposit(DTO *models.DepositDTO) (string, error)
	Create(DTO *models.DepositDTO) error
}

type Server struct {
	protobuff.PayServer
	service Uniplemented
}

func Register(gRPC *grpc.Server, service Uniplemented) {
	protobuff.RegisterPayServer(gRPC, &Server{service: service})
}

func NewServer(service Uniplemented) *Server {
	return &Server{
		service: service,
	}
}
