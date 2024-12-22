package postgresql

import (
	"CW_DB_v2/internal/domain/models"
	//"CW_DB_v2/internal/storage"
	"context"
	//"database/sql"
	//"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(DBconnection string) (*Storage, error) {

	const op = "storage.psql.New"

	dbPool, err := pgxpool.New(context.Background(), DBconnection)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Release()

	err = conn.Conn().Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Println("Успешное подключение к базе данных.")
	return &Storage{db: dbPool}, nil
}

func (s *Storage) RunMigrations() error {
	const op = "storage.psql.RunMigrations"

	migrationsPath := "migrations"
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(migrationsPath, file.Name())
			if err := executeSQLFile(s.db, filePath); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			log.Printf("Миграция %s успешно выполнена", file.Name())
		}
	}

	return nil

}

func executeSQLFile(dbPool *pgxpool.Pool, filePath string) error {

	const op = "storage.psql.RunMigrations.executeSQLFile"

	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	ctx := context.Background()
	_, err = dbPool.Exec(ctx, string(sqlContent))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveUser(ctx context.Context, login string, passHash []byte) (int64, error) {
	const op = "storage.postgresql.SaveUser"

	query := "SELECT register_client(@user_name, @pass_hash)"
	args := pgx.NamedArgs{
		"user_name": login,
		"pass_hash": passHash,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	//id, err := res.lastIndex()
	//if err != nil {
	//	return 0, fmt.Errorf("%s : %w", op, err)
	//}

	return 52, nil

}

func (s *Storage) User(ctx context.Context, login string) (models.User, error) {
	const op = "storage.postgresql.User"

	query := "SELECT id, login, password_hash From Users WHERE login = $1"

	rows, err := s.db.Query(ctx, query, login)
	if err != nil {
		return models.User{}, fmt.Errorf("%s : %w", op, err)
	}

	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Login, &user.Pass_hash)
		if err != nil {
			return models.User{}, fmt.Errorf("%s : %w", op, err)
		}
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	return false, nil

}
