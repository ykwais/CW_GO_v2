package main

import (
	"CW_DB_v2/internal/app"
	"CW_DB_v2/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("start app", slog.Any("cfg", cfg), slog.Int("port", cfg.GRPC.Port))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	//hgjgjgh
	go application.GRPCServer.MustRun()

	//TODO: init app
	//TODO: init server

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping app", slog.String("signal", sign.String()))

	application.GRPCServer.Stop()

	log.Info("application stopped")
}

func setupLogger() *slog.Logger {

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

}
