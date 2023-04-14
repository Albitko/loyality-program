package controller

import (
	"context"

	"github.com/gin-gonic/gin"
)

type userAuthenticator interface {
	CheckIsLoginFree(ctx context.Context) error
	Register(ctx context.Context) error
	Auth(ctx context.Context) error
}

type userAuthHandler struct {
	auth userAuthenticator
}

func (u *userAuthHandler) checkRequestFormat() error {
	return nil
}

func (u *userAuthHandler) Register(c *gin.Context) {
	err := u.checkRequestFormat()
	if err != nil {
		return
	}
}

func (u *userAuthHandler) Login(c *gin.Context) {
	err := u.checkRequestFormat()
	if err != nil {
		return
	}
}

func NewUserAuthHandler(auth userAuthenticator) *userAuthHandler {
	return &userAuthHandler{
		auth: auth,
	}
}
