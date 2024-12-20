package app

import (
	grpcapp "CW_DB_v2/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTl time.Duration) *App {
	//TODO:init storage

	//TODO init cw service

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{GRPCServer: grpcApp}
}
