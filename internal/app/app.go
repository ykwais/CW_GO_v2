package app

import (
	grpcapp "CW_DB_v2/internal/app/grpc"
	"CW_DB_v2/internal/services/cw"
	"CW_DB_v2/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTl time.Duration) *App {

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	cwService := cw.New(log, storage, storage, storage, tokenTTl)

	//TODO init cw service

	grpcApp := grpcapp.New(log, cwService, grpcPort)

	return &App{GRPCServer: grpcApp}
}
