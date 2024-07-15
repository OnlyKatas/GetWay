package controller

import (
	"bank/internal/models"
	protobuff "bank/pkg/proto/bank/gen/go"
	"context"
	"google.golang.org/grpc"
)

type Plug interface {
	ReqForPayProcessing(req *models.Request) (string, error)
	ReqForPayStatus(req *models.Request) (int, error)
}

type Server struct {
	protobuff.BankServer
	service Plug
}

func Register(gRPC *grpc.Server, service Plug) {
	protobuff.RegisterBankServer(gRPC, &Server{service: service})
}

func NewServer(service Plug) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) RequestForPaymentProcessing(ctx context.Context, req *protobuff.RFPPRequest) (resp *protobuff.RFPPResponse, err error) {
	request := models.Request{
		UserUUID: req.GetUserUuid(),
		Amount:   int(req.GetAmount()),
	}
	bankUUID, err := s.service.ReqForPayProcessing(&request)
	if err != nil {
		return nil, err
	}

	resp = &protobuff.RFPPResponse{
		BankUuid: bankUUID,
	}
	return resp, nil
}

func (s *Server) RequestForPayStatus(ctx context.Context, req *protobuff.RFPSRequest) (resp *protobuff.RFPSResponse, err error) {
	request := models.Request{
		UserUUID: req.GetUserUuid(),
		BankUUID: req.GetBankUuid(),
	}
	status, err := s.service.ReqForPayStatus(&request)
	if err != nil {
		return nil, err
	}
	resp = &protobuff.RFPSResponse{
		Status: int64(status),
	}
	return resp, nil
}
