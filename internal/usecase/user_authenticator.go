package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

type userRepository interface {
	Register(ctx context.Context, id, login, hashedPassword string) error
	GetCredentials(ctx context.Context, login string) (string, error)
}

type authenticator struct {
	repository userRepository
}

func (a *authenticator) CheckIsLoginFree(ctx context.Context, login string) error {
	_, err := a.repository.GetCredentials(ctx, login)
	return err
}

func (a *authenticator) Register(ctx context.Context, login, password string) error {
	id := uuid.New().String()
	hash := sha256.New()
	hash.Write([]byte(password))
	err := a.repository.Register(ctx, id, login, hex.EncodeToString(hash.Sum(nil)))
	return err
}

func (a *authenticator) Auth(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthenticator(repository userRepository) *authenticator {
	return &authenticator{
		repository: repository,
	}
}
