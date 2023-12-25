package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/Len4i/auth-service/internal/app/grpc"
	"github.com/Len4i/auth-service/internal/services/auth"
	"github.com/Len4i/auth-service/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		return nil
	}

	authSvc := auth.NewAuth(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcApp.NewApp(log, port, authSvc)
	return &App{
		GRPCApp: grpcApp,
	}
}
