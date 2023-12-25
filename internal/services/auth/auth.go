package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Len4i/auth-service/internal/domain/models"
	"github.com/Len4i/auth-service/internal/lib/jwt"
	"github.com/Len4i/auth-service/internal/services/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorInvalidAppID       = errors.New("invalid app id")
	ErrorInvalidUserID      = errors.New("invalid user id")
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (app models.App, err error)
}

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

// NewAuth creates new auth service
func NewAuth(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// RegisterNewUser registers new user
//
// If user with such email already exists, returns error
func (a *Auth) Register(ctx context.Context, email string, password string) (userID int64, err error) {
	const op = "auth.Register"
	log := a.log.With(slog.String("operation", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err = a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrorUserExists) {
			log.Warn("user already exists", slog.String("email", email))
			return 0, fmt.Errorf("%s: %w", op, ErrorInvalidUserID)
		}
		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered", slog.Int64("userID", userID))

	return userID, nil

}

// Login logs user in and returns token
//
// If user is not found or password is incorrect, returns error
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (token string, err error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("operation", op))

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			a.log.Warn("user not found", slog.String("email", email))
			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("failed to compare password", err)
		return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrorAppNotFound) {
			a.log.Warn("app not found", slog.Int("appID", appID))
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}
	log.Info("user logged in", slog.Int64("userID", user.ID))

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// IsAdmin checks if user is admin
//
// If user is not found, returns error
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(slog.String("operation", op))

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorAppNotFound) {
			a.log.Warn("user not found", slog.Int64("userID", userID))
			return false, fmt.Errorf("%s: %w", op, ErrorInvalidAppID)
		}
		log.Error("failed to get user", err)
		return false, err
	}

	return isAdmin, nil
}
