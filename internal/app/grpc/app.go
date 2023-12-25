package grpcApp

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	authgRPC "github.com/Len4i/auth-service/internal/grpc/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int, authSvc authgRPC.Auth) *App {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()))
	authgRPC.Register(grpcServer, authSvc)
	return &App{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	const op = "grpcApp.Run"
	log := a.log.With(slog.String("operation", op))

	log.Info("starting grpc server", slog.Int("port", a.port))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error("failed to tcp listener server", err)
		os.Exit(1)
	}
	log.Info("grpc server started", slog.String("address", l.Addr().String()))
	if err := a.grpcServer.Serve(l); err != nil {
		log.Error("failed to start grpc server", err)
		os.Exit(1)
	}
}

func (a *App) Stop() error {
	const op = "grpcApp.Stop"
	log := a.log.With(slog.String("operation", op))

	log.Info("stopping grpc server")
	a.grpcServer.GracefulStop()
	log.Info("grpc server stopped")

	return nil
}
