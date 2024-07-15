package run

import (
	"auth/internal/modules/gRPC/auth"
	"auth/internal/modules/gRPC/user"
	"auth/internal/modules/service"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int, addr string) *App {
	gRPCServer := grpc.NewServer()
	UserClient, err := user.NewUserClient(addr)
	if err != nil {
		log.Info("Failed to create User Client", err)
	}
	authSev := service.NewAuthService(UserClient)

	auth.Register(gRPCServer, authSev)
	auth.NewServerAPI(authSev)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}

}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("operation", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("operation", op)).
		Info("grpc server is stopping", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
