package main

import (
	"getway/config"
	"getway/run"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger()

	app := run.NewServer()
	if err := app.Run(cfg); err != nil {
		log.Info("Failed to start server", err)
	}
	go func() {
		if err := app.Run(cfg); err != nil {
			log.Info("Failed to start server", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("shutting down", slog.Any("signal", sign))

	app.Stop()
	log.Info("app stopped")

}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
