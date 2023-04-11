package usecase

import (
	"context"
)

type userRepository interface {
	CheckLoginRegister(ctx context.Context, login string) error
	Register(ctx context.Context, login, password string) error
	GetCredentials(ctx context.Context, login string) (string, error)
}

type authenticator struct {
	repository userRepository
}

func (a *authenticator) CheckIsLoginFree(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a *authenticator) Register(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
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
