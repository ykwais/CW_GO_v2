package cw

import (
	"CW_DB_v2/internal/domain/models"
	"CW_DB_v2/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type CW struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	//appProvider AppProvider
	tokenTTl time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, login string) (models.User, error)
	isAdmin(ctx context.Context, userID int64) (bool, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

//type AppProvider interface {
//	App(ctx context.Context, appId int) (models.App, error)
//}

func New(log *slog.Logger, userSaver UserSaver, usrProvider UserProvider /* appProvider AppProvider,*/, tokenTTl time.Duration) *CW {
	return &CW{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: usrProvider,
		//appProvider: appProvider,
		tokenTTl: tokenTTl,
	}
}

func (cw *CW) Login(ctx context.Context, login string, password string) (string, error) {
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
		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	log.Info("user registered successfully")

	return id, nil

}

func (cw *CW) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("implement me")
}
