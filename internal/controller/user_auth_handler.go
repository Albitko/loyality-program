package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Albitko/loyalty-program/internal/entities"
)

type userAuthenticator interface {
	CheckIsLoginFree(ctx context.Context, login string) error
	Register(ctx context.Context, login, password string) error
	Auth(ctx context.Context) error
}

type userAuthHandler struct {
	auth userAuthenticator
}

func (u *userAuthHandler) Register(c *gin.Context) {
	var request entities.AuthRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{Message: err.Error()})
		return
	}

	err = u.auth.Register(c, request.Login, request.Password)
	if errors.Is(err, entities.ErrLoginAlreadyInUse) {
		c.JSON(http.StatusConflict, entities.ErrorResponse{Message: "User already exists with the given login"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ErrorResponse{Message: "User registered"})
}

func (u *userAuthHandler) Login(c *gin.Context) {

}

func NewUserAuthHandler(auth userAuthenticator) *userAuthHandler {
	return &userAuthHandler{
		auth: auth,
	}
}
