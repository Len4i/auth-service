package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Len4i/auth-service/internal/app"
	"github.com/Len4i/auth-service/internal/config"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	// Init logger
	// httpLog used in middleware for http access logs
	log := setupLogger(cfg.Env)

	log.Info(
		"starting versions-collector",
		slog.String("env", cfg.Env),
		slog.String("version", "0.0.1"),
	)
	log.Debug("debug messages are enabled")

	application := app.NewApp(log, cfg.StoragePath, cfg.TokenTTL, cfg.GRPC.Port)
	go application.GRPCApp.MustRun()

	// Channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	_, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := application.GRPCApp.Stop(); err != nil {
		log.Error("failed to stop grpc server", err)
		os.Exit(1)
	}

	log.Info("exiting the app")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			}),
		)

	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log
}
