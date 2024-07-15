package controller

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"session/internal/models"

	protobuff "session/pkg/proto/session/gen/go"
)

type Sessioner interface {
	SendSessionToBank(ses *models.Session) (string, error)
}

type Server struct {
	protobuff.SessionServer
	service Sessioner
}

func Register(gRPC *grpc.Server, service Sessioner) {
	protobuff.RegisterSessionServer(gRPC, &Server{service: service})
}

func NewServer(service Sessioner) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) SendSession(ctx context.Context, req *protobuff.SendSessionRequest) (res *protobuff.SendSessionResponse, err error) {
	session := models.Session{
		ID:       int(req.GetId()),
		UserID:   int(req.GetUserId()),
		UserUUID: req.GetUserUuid(),
		Amount:   int(req.GetAmount()),
	}
	resp, err := s.service.SendSessionToBank(&session)
	if err != nil {
		return nil, fmt.Errorf("error send session to bank: %w", err)
	}
	res = &protobuff.SendSessionResponse{
		Req: resp,
	}
	return res, err
}
