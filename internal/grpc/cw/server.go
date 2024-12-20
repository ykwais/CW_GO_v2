package cw

import (
	"context"
	"fmt"
	cwv1 "github.com/ykwais/CW_GO_protos/gen/go/cw"
	"google.golang.org/grpc"
)

type serverAPI struct {
	cwv1.UnimplementedServiceServer
}

func RegisterServerAPI(gRPC *grpc.Server) {
	cwv1.RegisterServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *cwv1.LoginRequest) (*cwv1.LoginResponse, error) {
	fmt.Println("login: " + req.Login)
	fmt.Println("password: " + req.Password)
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *cwv1.RegisterRequest) (*cwv1.RegisterResponse, error) {
	fmt.Println("login: " + req.Login)
	fmt.Println("password: " + req.Password)
	//panic("implement me")
	return &cwv1.RegisterResponse{
		UserId: 123,
	}, nil
}

func (s *serverAPI) isAdmin(ctx context.Context, req *cwv1.IsAdminRequest) (*cwv1.IsAdminResponse, error) {
	panic("implement me")
}
