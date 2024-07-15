package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"session/config"
	"session/run"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger()

	application := run.NewApp(log, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	go application.MustRun(ctx)

	<-ctx.Done()

	log.Info("shutting down", slog.Any("signal", ctx.Err()))

	application.Stop()
	log.Info("app stopped")
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
