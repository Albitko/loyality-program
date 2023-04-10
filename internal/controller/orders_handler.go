package controller

import (
	"github.com/gin-gonic/gin"
)

type ordersHandler struct{}

func (o *ordersHandler) CreateOrder(c *gin.Context) {
}

func (o *ordersHandler) GetOrders(c *gin.Context) {
}

func (o *ordersHandler) GetBalance(c *gin.Context) {
}

func (o *ordersHandler) Withdraw(c *gin.Context) {
}

func (o *ordersHandler) GetWithdrawn(c *gin.Context) {
}

func NewOrdersHandler() *ordersHandler {
	return &ordersHandler{}
}
