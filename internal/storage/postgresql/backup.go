package postgresql

import (
	"fmt"
	"os"
	"os/exec"
)

func BackupDatabase(dbName, outputFile, user, password string) error {
	cmd := exec.Command("pg_dump", "-U", user, "-h", "localhost", "-p", "5432", "-d", dbName, "-f", outputFile)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка создания резервной копии: %w", err)
	}
	return nil
}

func RestoreDatabase(dbName, inputFile, user, password string) error {
	cmd := exec.Command("psql", "-U", user, "-h", "localhost", "-p", "5432", "-d", dbName, "-f", inputFile)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка восстановления базы данных: %w", err)
	}
	return nil
}
