package run

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	bank "session/BANK/Plug"
	"session/config"
	"session/internal/modules/service"
	"session/internal/modules/storage"
	"time"

	"google.golang.org/grpc"
	"net"
	controller "session/internal/modules/gRPC"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, cfg *config.Config, db *sqlx.DB) *App {
	gRPCServer := grpc.NewServer()
	Bank := bank.NewBANK()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if err := Bank.CanceledOP(); err != nil {
				return
			}
		}
	}()
	Storage := storage.NewBankFacade(Bank)
	Service := service.NewService(Storage)

	controller.Register(gRPCServer, Service)
	controller.NewServer(Service)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfg.Local.Port,
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
