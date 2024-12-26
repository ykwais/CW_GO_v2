package postgresql

import (
	"CW_DB_v2/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"strconv"
	"time"

	//"github.com/jackc/pgx/v5"

	//"CW_DB_v2/internal/storage"
	"context"
	//"database/sql"
	//"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	db *pgxpool.Pool
}

type CanMoney struct {
	Data      []byte
	VehicleId int64
	Brand     string
	Model     string
	TotalCost string
	Url       string
}

func New(DbConnection string) (*Storage, error) {

	const op = "storage.psql.New"

	dbPool, err := pgxpool.New(context.Background(), DbConnection)
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

func (s *Storage) GetUserBookings(userId int64) ([]models.UserBooking, error) {
	const op = "storage.psql.GetUserBookings"

	query := "Select * from get_user_bookings(@user_id);"

	args := pgx.NamedArgs{"user_id": userId}

	rows, err := s.db.Query(context.Background(), query, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var res []models.UserBooking
	for rows.Next() {
		var current models.UserBooking
		var start time.Time
		var end time.Time
		err := rows.Scan(&current.VehicleID, &current.Brand, &current.Model, &start, &end)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		current.StartDate = start.Format("2006-01-02")
		current.EndDate = end.Format("2006-01-02")

		res = append(res, current)
	}
	return res, nil

}

func (s *Storage) SelectAuto(userId int64, vehicleId int64, dateStart string, dateEnd string) (int64, error) {
	const op = "storage.psql.SelectAuto"
	strId := strconv.FormatInt(userId, 10)
	strEh := strconv.FormatInt(vehicleId, 10)

	fmt.Println("!!!!!!!!!!!! " + strId + " " + strEh + " " + dateStart + " " + dateEnd)

	//query := "Select book_vehicle(@p_user_id, @p_vehicle_id, @p_date_begin, @p_date_end)"
	//args := pgx.NamedArgs{
	//	"@p_user_id":    strId,
	//	"@p_vehicle_id": strEh,
	//	"@p_date_begin": dateStart,
	//	"@p_date_end":   dateEnd,
	//}

	query := "Select book_vehicle($1, $2, $3, $4)"
	args := []interface{}{
		userId,
		vehicleId,
		dateStart,
		dateEnd,
	}

	var bookingId int64
	err := s.db.QueryRow(context.Background(), query, args...).Scan(&bookingId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	//for rows.Next() {
	//	err := rows.Scan(&bookingId)
	//	if err != nil {
	//		return 0, err
	//	}
	//}

	fmt.Println(bookingId)

	return bookingId, nil

}

func (s *Storage) SaveUser(ctx context.Context, login string, passHash []byte, email string, real_name string) (int64, error) {
	const op = "storage.postgresql.SaveUser"

	resultCh := make(chan struct {
		id  int64
		err error
	}, 1)

	go func() {
		query := "SELECT register_user(@user_name, @pass_hash, @email, @real_name)"
		args := pgx.NamedArgs{
			"user_name": login,
			"pass_hash": passHash,
			"email":     email,
			"real_name": real_name,
		}

		var id_user int64
		err := s.db.QueryRow(ctx, query, args).Scan(&id_user)
		if err != nil {
			resultCh <- struct {
				id  int64
				err error
			}{0, fmt.Errorf("%s : %w", op, err)}
			return
		}

		resultCh <- struct {
			id  int64
			err error
		}{id_user, nil}
	}()

	result := <-resultCh

	return result.id, result.err

	//query := "SELECT register_client(@user_name, @pass_hash)"
	//args := pgx.NamedArgs{
	//	"user_name": login,
	//	"pass_hash": passHash,
	//}
	//
	//_, err := s.db.Exec(ctx, query, args)
	//if err != nil {
	//	return 0, fmt.Errorf("%s : %w", op, err)
	//}
	//
	////id, err := res.lastIndex()
	////if err != nil {
	////	return 0, fmt.Errorf("%s : %w", op, err)
	////}
	//
	//return 52, nil

}

func (s *Storage) PhotosOfOneAutomobile(id int64) ([]models.Photo, error) {
	const op = "storage.photosOfOneAutomobile"

	query := "Select * from get_vehicle_photos_table(@id)"
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := s.db.Query(context.Background(), query, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var res []models.Photo
	for rows.Next() {
		var result models.Photo
		err := rows.Scan(&result.Name) //тут просто путь сохраняю
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, result)
	}

	return res, nil
}

func (s *Storage) GetAvailableCars(start_time string, end_time string) ([]models.BetterPhoto, error) {
	const op = "storage.psql.GetAvailableCars"

	query := "Select * from get_available_vehicles(@start, @end)"
	args := pgx.NamedArgs{
		"start": start_time,
		"end":   end_time,
	}

	rows, err := s.db.Query(context.Background(), query, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var res []models.BetterPhoto
	for rows.Next() {
		var result CanMoney
		err := rows.Scan(&result.VehicleId, &result.Brand, &result.Model, &result.TotalCost, &result.Url)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		costStr := result.TotalCost
		costStr = strings.TrimPrefix(costStr, "$")
		totalCost, err := strconv.ParseFloat(costStr, 64) // Преобразуем в float64
		if err != nil {
			log.Println("failed to convert price to float", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		res = append(res, models.BetterPhoto{
			VehicleId: result.VehicleId,
			Brand:     result.Brand,
			Model:     result.Model,
			TotalCost: totalCost,
			Url:       result.Url,
		})
	}
	return res, nil
}

func (s *Storage) User(ctx context.Context, login string) (models.User, error) {
	const op = "storage.postgresql.User"

	resultCh := make(chan struct {
		user models.User
		err  error
	}, 1)

	go func() {

		query := "SELECT * from login_user(@login)"
		args := pgx.NamedArgs{
			"login": login,
		}

		rows, err := s.db.Query(ctx, query, args)
		if err != nil {
			resultCh <- struct {
				user models.User
				err  error
			}{models.User{}, fmt.Errorf("%s : %w", op, err)}
			return
		}

		var user models.User

		for rows.Next() {
			err := rows.Scan(&user.ID, &user.Login, &user.Pass_hash)
			if err != nil {

				resultCh <- struct {
					user models.User
					err  error
				}{models.User{}, fmt.Errorf("%s : %w", op, err)}
				return
			}
		}

		resultCh <- struct {
			user models.User
			err  error
		}{user, nil}
	}()

	result := <-resultCh

	return result.user, result.err

	//query := "SELECT * from login_user(@login)"
	//args := pgx.NamedArgs{
	//	"login": login,
	//}
	//
	//rows, err := s.db.Query(ctx, query, args)
	//if err != nil {
	//	return models.User{}, fmt.Errorf("%s : %w", op, err)
	//}
	//
	//var user models.User
	//for rows.Next() {
	//	err := rows.Scan(&user.ID, &user.Login, &user.Pass_hash)
	//	if err != nil {
	//		return models.User{}, fmt.Errorf("%s : %w", op, err)
	//	}
	//}
	//
	//return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	return false, nil

}
