package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stasdashkevitch/wizzard/auth/internal/entity"
	"github.com/stasdashkevitch/wizzard/auth/internal/lib/jwt"
	"github.com/stasdashkevitch/wizzard/auth/internal/repository"
	"github.com/stasdashkevitch/wizzard/common/logger"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          logger.ILogger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (entity.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (entity.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app_id")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new instance of the Auth service.
func New(
	log logger.ILogger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exists, returns error.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	const op = "auth.Login"

	a.log.Info("attempting to login user", "op", op, "email", email)

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.log.Warn("user not found", "error", err.Error(), "op", op, "email", email)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", "error", err.Error(), "op", op, "email", email)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Error("invalid credentials", "error", err.Error(), "op", op, "email", email)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		a.log.Error("failed to get app_id", "error", err.Error(), "op", op, "email", email)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", "error", err.Error(), "op", op, "email", email)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("user logged in successfully", "op", op, "email", email)

	return token, nil
}

// RegisterNewUser registers new user in the system and returns new ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	a.log.Info("registering user", "op", op, "email", email)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to generate passsword hash", "error", err.Error(), "op", op, "email", email)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			a.log.Warn("user already exists", "error", err, "op", op, "email", email)

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		a.log.Error("failed to save user", "error", err.Error(), "op", op, "email", email)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("user registered", "op", op, "email", email)

	return id, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	a.log.Info("checking if user is admin", "op", op, "user_id", userID)

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			a.log.Warn("app not found", "error", err.Error(), "op", op, "user_id", userID)
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		a.log.Error("failed to check if user is admin", "error", err.Error(), "op", op, "user_id", userID)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
