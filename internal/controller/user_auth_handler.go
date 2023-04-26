package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Albitko/loyalty-program/internal/entities"
	"github.com/Albitko/loyalty-program/internal/utils"
)

type userAuthenticator interface {
	Register(ctx context.Context, user entities.User) error
	Auth(ctx context.Context, login, password string) (entities.User, error)
	CreateAccessToken(user entities.User) (string, error)
}

type userAuthHandler struct {
	auth userAuthenticator
}

func (u *userAuthHandler) Register(c *gin.Context) {
	var request entities.AuthRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		utils.Logger.Error("userAuthHandler:Register - request bind JSON ", zap.Error(err))
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{Message: err.Error()})
		return
	}

	user := entities.User{
		ID:       uuid.New().String(),
		Login:    request.Login,
		Password: request.Password,
	}

	err = u.auth.Register(c, user)
	if errors.Is(err, entities.ErrLoginAlreadyInUse) {
		utils.Logger.Error("userAuthHandler:Register - login already in use", zap.Error(err))
		c.JSON(http.StatusConflict, entities.ErrorResponse{Message: "User already exists with the given login"})
		return
	}
	if err != nil {
		utils.Logger.Error("userAuthHandler:Register - Register ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, err := u.auth.CreateAccessToken(user)
	if err != nil {
		utils.Logger.Error("userAuthHandler:Register - CreateAccessToken ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.Header("Authorization", accessToken)
	c.JSON(http.StatusOK, entities.ErrorResponse{Message: "User registered"})
}

func (u *userAuthHandler) Login(c *gin.Context) {
	var request entities.AuthRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		utils.Logger.Error("userAuthHandler:Login - request bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := u.auth.Auth(c, request.Login, request.Password)
	if errors.Is(err, entities.ErrInvalidCredentials) {
		utils.Logger.Error("userAuthHandler:Login - wrong credentials", zap.Error(err))
		c.JSON(http.StatusUnauthorized, entities.ErrorResponse{Message: "Invalid login or password"})
		return
	}
	if err != nil {
		utils.Logger.Error("userAuthHandler:Login - Auth", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, err := u.auth.CreateAccessToken(user)
	if err != nil {
		utils.Logger.Error("userAuthHandler:Login - CreateAccessToken", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.Header("Authorization", accessToken)
	c.JSON(http.StatusOK, entities.ErrorResponse{Message: "User registered"})
}

func NewUserAuthHandler(auth userAuthenticator) *userAuthHandler {
	return &userAuthHandler{
		auth: auth,
	}
}
