package main

import (
	"log/slog"
	"os/signal"
	"payS/config"
	"payS/internal/infrastructure/db/postgres"
	"payS/run"
	"syscall"

	"os"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Info("Failed to connect to database", err)
	}

	application := run.NewApp(log, cfg, db)

	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("shutting down", slog.Any("signal", sign))

	application.Stop()
	log.Info("app stopped")

}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
