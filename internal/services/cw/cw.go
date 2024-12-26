package cw

import (
	"CW_DB_v2/internal/domain/models"
	"CW_DB_v2/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"path/filepath"
)

type CW struct {
	log  *slog.Logger
	srvc Service
}

type Service interface {
	SaveUser(ctx context.Context, login string, passHash []byte, email string, real_name string) (uid int64, err error)
	User(ctx context.Context, login string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	GetAvailableCars(start_time string, end_time string) ([]models.BetterPhoto, error)
	PhotosOfOneAutomobile(id int64) ([]models.Photo, error)
	SelectAuto(userId int64, vehicleId int64, dateStart string, dateEnd string) (bookingId int64, err error)
	GetUserBookings(userId int64) ([]models.UserBooking, error)
	CancelBooking(userId int64, vehicleId int64) (bool, error)
	GetDataForAdmin() ([]models.AdminData, error)
	GetUsersForAdmin() ([]models.BetterUser, error)
	DeleteUser(id int64) (bool, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

func New(log *slog.Logger, service Service) *CW {
	return &CW{
		log:  log,
		srvc: service,
	}
}

/*
ниже представлены уже сама реализация обработки запроса, то есть мы получаем входные данные из реквеста и перенаправляем их в сущность, которая взаимодействует с бд
*/

func (cw *CW) DeleteUser(userId int64) (bool, error) {
	cw.log.Info("DeleteUser")
	res, err := cw.srvc.DeleteUser(userId)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (cw *CW) CancelBooking(userId int64, vehicleId int64) (bool, error) {
	cw.log.Info("cancel booking start")
	res, err := cw.srvc.CancelBooking(userId, vehicleId)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (cw *CW) GetUsersForAdmin() ([]models.BetterUser, error) {
	cw.log.Info("get users for admin start")

	users, err := cw.srvc.GetUsersForAdmin()
	if err != nil {
		return nil, err
	}
	return users, nil

}

func (cw *CW) GetDataForAdmin() ([]models.AdminData, error) {
	cw.log.Info("get data for admin start")
	infos, err := cw.srvc.GetDataForAdmin()
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (cw *CW) GetUserBookings(userId int64) ([]models.UserBooking, error) {
	cw.log.Info("get user bookings for userId=%d", userId)
	bookings, err := cw.srvc.GetUserBookings(userId)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (cw *CW) PhotosOfAutomobile(id int64) (photos []models.Photo, err error) {
	cw.log.Info("starting inner photos of autos")
	returnedPhotos, err := cw.srvc.PhotosOfOneAutomobile(id)
	if err != nil {
		return nil, err
	}

	for _, photo := range returnedPhotos {
		if _, err := os.Stat(photo.Name); os.IsNotExist(err) {
			return nil, err
		}
		data, err := os.ReadFile(photo.Name)
		if err != nil {
			return nil, err
		}

		photos = append(photos, models.Photo{
			Name: photo.Name, //тут путь
			Data: data,
		})
	}

	return photos, err
}

func (cw *CW) PhotosForMainScreen(ctx context.Context, data_start string, data_end string) (photos []models.BetterPhoto, err error) {
	cw.log.Info("starting photos for main screen")
	better_photos, err := cw.srvc.GetAvailableCars(data_start, data_end)
	if err != nil {
		return nil, err
	}

	for _, better_photo := range better_photos {
		if _, err := os.Stat(better_photo.Url); os.IsNotExist(err) {
			return nil, err
		}

		data, err := os.ReadFile(better_photo.Url)
		if err != nil {
			return nil, err
		}

		photos = append(photos, models.BetterPhoto{
			Data:      data,
			Url:       better_photo.Url,
			VehicleId: better_photo.VehicleId,
			Model:     better_photo.Model,
			Brand:     better_photo.Brand,
			TotalCost: better_photo.TotalCost,
		})

	}

	return photos, nil

}

func (cw *CW) ListPhotos() ([]models.Photo, error) {

	cw.log.Info("starting list photos")

	photosDir := "./photos"

	if _, err := os.Stat(photosDir); os.IsNotExist(err) {
		return nil, err
	}

	files, err := os.ReadDir(photosDir)
	if err != nil {
		return nil, err
	}

	var photos []models.Photo
	for _, file := range files {

		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(photosDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		photos = append(photos, models.Photo{
			Name: file.Name(),
			Data: data,
		})
	}

	return photos, nil

}

func (cw *CW) Login(ctx context.Context, login string, password string) (int64, error) {
	const op = "cw.Login"

	log := cw.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("logining user")

	user, err := cw.srvc.User(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			cw.log.Warn("user not found", slog.String("login", login), slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)

		}

		cw.log.Error("failed to get user", slog.String("login", login), slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)

	}

	if err := bcrypt.CompareHashAndPassword(user.Pass_hash, []byte(password)); err != nil {
		cw.log.Info("invalid credentials", slog.String("login", login), slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully", slog.String("login", login))

	//token, err := jwt.NewToken(user /* app, */)
	//if err != nil {
	//	cw.log.Error("failed to create token", slog.String("login", login), slog.String("error", err.Error()))
	//	return "", fmt.Errorf("%s: %w", op, err)
	//}

	return user.ID, nil

}

func (cw *CW) SelectAuto(userId int64, vehicleId int64, dateStart string, dateEnd string) (bookingId int64, err error) {
	const op = "cw.SelectAuto"

	id, err := cw.srvc.SelectAuto(userId, vehicleId, dateStart, dateEnd)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (cw *CW) Register(ctx context.Context, login string, password string, email string, real_name string) (int64, error) {
	const op = "cw.RegisterNewUser"

	log := cw.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("registering new user")

	resultCh := make(chan struct {
		id  int64
		err error
	}, 1)

	go func() {

		passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate hash password", slog.String("error", err.Error()))
			resultCh <- struct {
				id  int64
				err error
			}{0, fmt.Errorf("%s : %w", op, err)}
			return
		}

		id, err := cw.srvc.SaveUser(ctx, login, passHash, email, real_name)
		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				log.Warn("user already exists", slog.String("login", login), slog.String("error", err.Error()))
				resultCh <- struct {
					id  int64
					err error
				}{0, fmt.Errorf("%s : %w", op, ErrUserExists)}
				return
			}
			log.Error("failed to save user", slog.String("error", err.Error()))
			resultCh <- struct {
				id  int64
				err error
			}{0, fmt.Errorf("%s : %w", op, err)}
			return
		}

		log.Info("user registered successfully")
		resultCh <- struct {
			id  int64
			err error
		}{id, nil}
	}()

	result := <-resultCh

	if result.err != nil {
		return 0, result.err
	}

	return result.id, nil

	//passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//if err != nil {
	//	log.Error("failed to generate hash password", slog.String("error", err.Error()))
	//
	//	return 0, fmt.Errorf("%s : %w", op, err)
	//}
	//
	//id, err := cw.usrSaver.SaveUser(ctx, login, passHash)
	//if err != nil {
	//	if errors.Is(err, storage.ErrUserExists) {
	//		log.Warn("user already exists", slog.String("login", login), slog.String("error", err.Error()))
	//
	//		return 0, fmt.Errorf("%s : %w", op, ErrUserExists)
	//	}
	//	log.Error("failed to save user", slog.String("error", err.Error()))
	//	return 0, fmt.Errorf("%s : %w", op, err)
	//}
	//
	//log.Info("user registered successfully")
	//
	//return id, nil
}

//func (cw *CW) RegisterNewUser(ctx context.Context, login string, password string) (int64, error) {
//	const op = "cw.RegisterNewUser"
//
//	log := cw.log.With(slog.String("op", op), slog.String("login", login))
//	log.Info("registering new user")
//
//	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		log.Error("failed to generate hash password", slog.String("error", err.Error()))
//
//		return 0, fmt.Errorf("%s : %w", op, err)
//	}
//
//	id, err := cw.usrSaver.SaveUser(ctx, login, passHash)
//	if err != nil {
//		if errors.Is(err, storage.ErrUserExists) {
//			log.Warn("user already exists", slog.String("login", login), slog.String("error", err.Error()))
//
//			return 0, fmt.Errorf("%s : %w", op, ErrUserExists)
//		}
//		log.Error("failed to save user", slog.String("error", err.Error()))
//		return 0, fmt.Errorf("%s : %w", op, err)
//	}
//
//	log.Info("user registered successfully")
//
//	return id, nil
//
//}

func (cw *CW) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "cw.IsAdmin"

	return false, nil

}
