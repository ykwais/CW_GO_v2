package cw

import (
	"context"
	cwv1 "github.com/ykwais/CW_GO_protos/gen/go/cw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type CW interface {
	Login(ctx context.Context, login, password string) (token string, err error)
	Register(ctx context.Context, login, password string) (userID int64, err error)
	isAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	cwv1.UnimplementedServiceServer
	*slog.Logger
	cw CW
}

func RegisterServerAPI(gRPC *grpc.Server, logger *slog.Logger, serv CW) {
	cwv1.RegisterServiceServer(gRPC, &serverAPI{
		UnimplementedServiceServer: cwv1.UnimplementedServiceServer{},
		Logger:                     logger,
		cw:                         serv,
	})

}

func (s *serverAPI) Login(ctx context.Context, req *cwv1.LoginRequest) (*cwv1.LoginResponse, error) {
	s.Logger.Info("on LOGIN request get: ", slog.String("login", req.Login), slog.String("password", req.Password))
	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty login or password")
	}

	token, err := s.cw.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		// TODO:...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &cwv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *cwv1.RegisterRequest) (*cwv1.RegisterResponse, error) {

	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty login or password")
	}

	userID, err := s.cw.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		//TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &cwv1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) isAdmin(ctx context.Context, req *cwv1.IsAdminRequest) (*cwv1.IsAdminResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be 0")
	}

	isAdmin, err := s.cw.isAdmin(ctx, req.GetUserId())
	if err != nil {
		//TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &cwv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
