package usecase

import (
	"context"
)

type authenticator struct {
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

func NewAuthenticator() *authenticator {
	return &authenticator{}
}
