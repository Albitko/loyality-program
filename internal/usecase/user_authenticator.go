package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/Albitko/loyalty-program/internal/entities"
	"github.com/Albitko/loyalty-program/internal/utils"
)

type userRepository interface {
	Register(ctx context.Context, id, login, hashedPassword string) error
	GetCredentials(ctx context.Context, login string) (entities.User, error)
}

type authenticator struct {
	repository userRepository
	secret     string
}

func (a *authenticator) CreateAccessToken(user entities.User) (string, error) {
	claims := &entities.JwtCustomClaims{
		Name: user.Login,
		ID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := unsignedToken.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (a *authenticator) Register(ctx context.Context, user entities.User) error {
	err := a.repository.Register(ctx, user.ID, user.Login, utils.HexHash(user.Password))
	return err
}

func (a *authenticator) Auth(ctx context.Context, login, password string) (entities.User, error) {

	passwordHash := utils.HexHash(password)
	user, err := a.repository.GetCredentials(ctx, login)
	if err != nil {
		return user, err
	}
	if passwordHash != user.Password {
		return user, entities.ErrInvalidCredentials
	}
	return user, nil
}

func NewAuthenticator(repository userRepository, secret string) *authenticator {
	return &authenticator{
		repository: repository,
		secret:     secret,
	}
}
