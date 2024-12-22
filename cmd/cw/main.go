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
	cfg := config.MustLoad() //читаем файл настройщик

	log := setupLogger() //подключаем логгер

	log.Info("start app", slog.Any("cfg", cfg), slog.Int("port", cfg.GRPC.Port))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL) //готовый сервер с логгером и сервисом

	go application.GRPCServer.MustRun() //если не смогли запуститься - паникуем

	stop := make(chan os.Signal, 1) //создаем канал слушателя сигналов системы
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop //ловим сигналы системы

	log.Info("stopping app", slog.String("signal", sign.String()))

	application.GRPCServer.Stop() //останавливаем сервер

	log.Info("application stopped")
}

func setupLogger() *slog.Logger {

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

}
