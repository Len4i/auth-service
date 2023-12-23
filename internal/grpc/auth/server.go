package auth

import (
	"context"

	aaav1 "github.com/Len4i/aaa/gen/go/aaa"
	"google.golang.org/grpc"
)

type ServerApi struct {
	aaav1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	aaav1.RegisterAuthServer(gRPC, &ServerApi{})
}

func (s *ServerApi) Login(ctx context.Context, req *aaav1.LoginRequest) (*aaav1.LoginResponse, error) {
	return &aaav1.LoginResponse{
		Token: "token",
	}, nil
}

func (s *ServerApi) Register(ctx context.Context, req *aaav1.RegisterRequest) (*aaav1.RegisterResponse, error) {
	panic("not implemented")
}

func (s *ServerApi) IsAdmin(ctx context.Context, req *aaav1.IsAdminRequest) (*aaav1.IsAdminResponse, error) {
	panic("not implemented")
}
