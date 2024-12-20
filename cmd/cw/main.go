package main

import (
	"CW_DB_v2/internal/app"
	"CW_DB_v2/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("start app", slog.Any("cfg", cfg), slog.Int("port", cfg.GRPC.Port))

	//TODO: init app

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GRPCServer.MustRun()

	//TODO: init server
}

func setupLogger() *slog.Logger {

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

}
