package grpcapp

import (
	cwgrpc "CW_DB_v2/internal/grpc/cw"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	cwgrpc.RegisterServerAPI(gRPCServer, log)

	return &App{log, gRPCServer, port}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(slog.String("op", op))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting grpc server", slog.String("addr", listener.Addr().String()))

	if err := app.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	log := app.log.With(slog.String("op", op))
	log.Info("stopping grpc server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()

}
