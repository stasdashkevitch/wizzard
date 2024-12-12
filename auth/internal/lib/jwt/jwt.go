package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stasdashkevitch/wizzard/auth/internal/entity"
)

func NewToken(user entity.User, app entity.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString(app.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// TODO: tests
