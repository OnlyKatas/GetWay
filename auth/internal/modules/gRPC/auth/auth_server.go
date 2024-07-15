package auth

import (
	global "auth/internal/models"
	protobuff "auth/pkg/proto/auth/gen/go"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserError struct {
	Message string
	Code    codes.Code
}

func (e *UserError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrUserNotFound       = UserError{"user not found", codes.NotFound}
	ErrInternal           = UserError{"internal server error", codes.Internal}
	ErrUserAlreadyExists  = UserError{"user already exists", codes.AlreadyExists}
	ErrInvalidCredentials = UserError{"invalid credentials", codes.Unauthenticated}
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	Register(ctx context.Context, user *global.User) (userID int64, uuid string, err error)
}

type ServerAPI struct {
	protobuff.AuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	protobuff.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

func NewServerAPI(auth Auth) *ServerAPI {
	return &ServerAPI{
		auth: auth,
	}
}

func (s *ServerAPI) Login(ctx context.Context, req *protobuff.LoginRequest) (*protobuff.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: ....Ошибка модет быть в том что введены ни коректнеые данные
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &protobuff.LoginResponse{Token: token}, nil

}

func (s *ServerAPI) Register(ctx context.Context, req *protobuff.RegisterRequest) (*protobuff.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	user := global.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Username:  req.GetUsername(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		UUID:      req.GetUuid(),
	}
	userID, uuid, err := s.auth.Register(ctx, &user)
	if err != nil {
		switch err.(type) {
		case *UserError:
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &protobuff.RegisterResponse{UserId: userID, Uuid: uuid}, nil

}

func validateLogin(req *protobuff.LoginRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "missing credentials")
	}
	return nil
}

func validateRegister(req *protobuff.RegisterRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "missing credentials")
	}
	return nil
}
