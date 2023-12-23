package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/Len4i/auth-service/internal/app/grpc"
)

type App struct {
	GRPCApp *grpcApp.App
}

func NewApp(
	log *slog.Logger,
	storagePath string,
	tokenTTL time.Duration,
	port int,
) *App {
	grpcApp := grpcApp.NewApp(log, port)
	return &App{
		GRPCApp: grpcApp,
	}
}
