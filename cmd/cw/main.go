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

	application := app.New(log, cfg.GRPC.Port, cfg.DbContainerPath) //готовый сервер с логгером и сервисом

	//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

	//dbName := "afdb"
	//user := "ykwais"
	//backupFile := "C:\\Users\\fedor\\GolandProjects\\CW_DB_v2\\storage\\backup.sql" // Путь на хосте
	//
	//// ID контейнера
	//containerID := "my_postgres_container" // Замените на ID вашего контейнера
	//
	//// Формируем команду для запуска pg_dump внутри контейнера
	//cmd := exec.Command(
	//	"docker", "exec", containerID, "pg_dump", "-U", user, "-F", "c", "-b", "-v", "-f", "/backup/backup.sql", dbName,
	//)
	//
	//// Перенаправляем вывод команды в консоль
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//
	//// Запускаем команду
	//err := cmd.Run()
	//if err != nil {
	//	fmt.Println("Ошибка при резервном копировании:", err)
	//	return
	//}
	//
	//// Копируем резервную копию с контейнера на хост
	//copyCmd := exec.Command("docker", "cp", containerID+":/backup/backup.sql", backupFile)
	//copyCmd.Stdout = os.Stdout
	//copyCmd.Stderr = os.Stderr
	//
	//// Выполняем копирование
	//err = copyCmd.Run()
	//if err != nil {
	//	fmt.Println("Ошибка при копировании файла резервной копии:", err)
	//	return
	//}
	//
	//fmt.Println("Резервная копия успешно создана:", backupFile)
	//
	////%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	//
	//dbName := "afdb"
	//user := "ykwais"
	//backupFile := "C:\\Users\\fedor\\GolandProjects\\CW_DB_v2\\storage\\backup.sql" // Путь к файлу резервной копии на хосте
	//
	//// ID контейнера
	//containerID := "my_postgres_container" // Замените на ID вашего контейнера
	//
	//// Проверяем, существует ли файл резервной копии
	//if _, err := os.Stat(backupFile); os.IsNotExist(err) {
	//	fmt.Println("Файл резервной копии не найден:", backupFile)
	//	return
	//}
	//
	//// Копируем файл резервной копии в контейнер
	//copyCmd := exec.Command("docker", "cp", backupFile, containerID+":/backup/backup.sql")
	//copyCmd.Stdout = os.Stdout
	//copyCmd.Stderr = os.Stderr
	//
	//err := copyCmd.Run()
	//if err != nil {
	//	fmt.Println("Ошибка при копировании файла резервной копии в контейнер:", err)
	//	return
	//}
	//
	//// Восстанавливаем базу данных из резервной копии с удалением существующих объектов
	//restoreCmd := exec.Command(
	//	"docker", "exec", containerID, "pg_restore", "-U", user, "-d", dbName, "--clean", "-v", "/backup/backup.sql",
	//)
	//
	//// Перенаправляем вывод команды в консоль
	//restoreCmd.Stdout = os.Stdout
	//restoreCmd.Stderr = os.Stderr
	//
	//// Запускаем восстановление
	//err = restoreCmd.Run()
	//if err != nil {
	//	fmt.Println("Ошибка при восстановлении базы данных:", err)
	//	return
	//}
	//
	//fmt.Println("База данных успешно восстановлена из резервной копии.")

	//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

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

//func BackupDatabase() error {
//	dbName := "afdb"           // Имя базы данных
//	outputFile := "backup.sql" // Файл, куда будет сохранена резервная копия
//
//	cmd := exec.Command("pg_dump", "-U", "ykwais", "-h", "localhost", "-p", "5432", "-d", dbName, "-f", outputFile)
//	cmd.Env = append(os.Environ(), "PGPASSWORD=1111")
//
//	if err := cmd.Run(); err != nil {
//		return fmt.Errorf("ошибка создания резервной копии: %w", err)
//	}
//
//	fmt.Println("Резервная копия успешно создана в файле:", outputFile)
//	return nil
//}
//
//func RestoreDatabase() error {
//	dbName := "afdb"          // Имя базы данных
//	inputFile := "backup.sql" // Файл с резервной копией
//
//	cmd := exec.Command("psql", "-U", "ykwais", "-h", "localhost", "-p", "5432", "-d", dbName, "-f", inputFile)
//	cmd.Env = append(os.Environ(), "PGPASSWORD=1111")
//
//	if err := cmd.Run(); err != nil {
//		return fmt.Errorf("ошибка восстановления базы данных: %w", err)
//	}
//
//	fmt.Println("Данные успешно восстановлены из файла:", inputFile)
//	return nil
//}
