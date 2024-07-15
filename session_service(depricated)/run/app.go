package run

import (
	"context"
	"log/slog"
	"session/config"
	payGRPC "session/internal/modules/gRPC/pay"
	"session/internal/modules/service"
)

type App struct {
	log        *slog.Logger
	payService *service.Service
}

func NewApp(log *slog.Logger, cfg *config.Config) *App {
	PayClient, err := payGRPC.NewClient(context.Background(), cfg.GRPC.PayAddr)
	if err != nil {
		log.Info("не могу создать клиента для платежного сервиса")
	}

	Service := service.NewService(PayClient, 10)

	return &App{
		log:        log,
		payService: Service,
	}
}

func (a *App) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "grpcapp.Run"
	a.payService.WorkerPool(ctx)
	a.log.Info("Воркер пул запущен")
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("operation", op)).
		Info("grpc server is stopping")
}
