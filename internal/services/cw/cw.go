package cw

import (
	"CW_DB_v2/internal/domain/models"
	"CW_DB_v2/internal/storage"
	"CW_DB_v2/lib/jwt"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type CW struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
}

type UserSaver interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

func New(log *slog.Logger, userSaver UserSaver, usrProvider UserProvider) *CW {
	return &CW{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: usrProvider,
	}
}

/*
	ниже представлены уже сама реализация обработки запроса, то есть мы получаем входные данные из реквеста и перенаправляем их в сущность, которая взаимодействует с бд
*/

func (cw *CW) Login(ctx context.Context, login string, password string /*, appID int*/) (string, error) {
	const op = "cw.Login"

	log := cw.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("logining user")

	user, err := cw.usrProvider.User(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			cw.log.Warn("user not found", slog.String("login", login), slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)

		}

		cw.log.Error("failed to get user", slog.String("login", login), slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)

	}

	if err := bcrypt.CompareHashAndPassword(user.Pass_hash, []byte(password)); err != nil {
		cw.log.Info("invalid credentials", slog.String("login", login), slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	/*app, err := cw.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}*/

	log.Info("user logged in successfully", slog.String("login", login))

	token, err := jwt.NewToken(user /* app, */)
	if err != nil {
		cw.log.Error("failed to create token", slog.String("login", login), slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}

func (cw *CW) Register(ctx context.Context, username string, password string) (int64, error) {
	return cw.RegisterNewUser(ctx, username, password)
}

func (cw *CW) RegisterNewUser(ctx context.Context, login string, password string) (int64, error) {
	const op = "cw.RegisterNewUser"

	log := cw.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash password", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s : %w", op, err)
	}

	id, err := cw.usrSaver.SaveUser(ctx, login, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("login", login), slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s : %w", op, ErrUserExists)
		}
		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	log.Info("user registered successfully")

	return id, nil

}

func (cw *CW) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "cw.IsAdmin"

	return false, nil

}
