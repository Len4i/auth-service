package auth

import (
	"context"

	aaav1 "github.com/Len4i/aaa/gen/go/aaa"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const emptyUserID = 0

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type ServerApi struct {
	aaav1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	aaav1.RegisterAuthServer(gRPC, &ServerApi{
		auth: auth,
	})
}

func (s *ServerApi) Login(ctx context.Context, req *aaav1.LoginRequest) (*aaav1.LoginResponse, error) {

	if err := validateRequestCreds(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: handle errors
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &aaav1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerApi) Register(ctx context.Context, req *aaav1.RegisterRequest) (*aaav1.RegisterResponse, error) {
	if err := validateRequestCreds(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &aaav1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *ServerApi) IsAdmin(ctx context.Context, req *aaav1.IsAdminRequest) (*aaav1.IsAdminResponse, error) {
	if err := validateIsAdmin(req.GetUserId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	ok, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &aaav1.IsAdminResponse{
		IsAdmin: ok,
	}, nil
}

func validateRequestCreds(email string, password string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateIsAdmin(userID int64) error {
	if userID == emptyUserID {
		return status.Error(codes.InvalidArgument, "userID is required")
	}
	return nil
}
