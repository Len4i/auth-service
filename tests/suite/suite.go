package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	aaav1 "github.com/Len4i/aaa/gen/go/aaa"
	"github.com/Len4i/auth-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	localHost = "localhost"
)

type Suite struct {
	*testing.T
	AuthClient aaav1.AuthClient
	Cfg        *config.Config
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../configs/local_tests.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(ctx, grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial grpc: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		AuthClient: aaav1.NewAuthClient(cc),
		Cfg:        cfg,
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(localHost, strconv.Itoa(cfg.GRPC.Port))
}
