package app

import "C"
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

	storage, err := sqlite.New(storagePath) //заменить на постгрю
	if err != nil {
		panic(err)
	}

	/*
		cwService - объект, в котором лежит логгер и наша сущность для общения с бд - то есть реализация общения самого - Методы Login, Register и тд
	*/
	cwService := cw.New(log, storage, storage, storage, tokenTTl) //тут много лишнего

	grpcApp := grpcapp.New(log, cwService, grpcPort) //итоговое приложение с логгером, сервером(с реализованным сервисом) и номером порта

	return &App{GRPCServer: grpcApp}
}
