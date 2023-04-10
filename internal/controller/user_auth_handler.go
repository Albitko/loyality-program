package controller

import (
	"github.com/gin-gonic/gin"
)

type userAuthHandler struct{}

func (u *userAuthHandler) Register(c *gin.Context) {
}

func (u *userAuthHandler) Login(c *gin.Context) {
}

func NewUserAuthHandler() *userAuthHandler {
	return &userAuthHandler{}
}
