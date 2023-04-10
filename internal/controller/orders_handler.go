package controller

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ordersProcessor interface {
	CheckOrderExist(ctx context.Context) error
	RegisterOrder(ctx context.Context) error
	GetUserOrder(ctx context.Context) error
}

type ordersHandler struct {
	processor ordersProcessor
}

func (o *ordersHandler) checkRequestFormat() error {
	return nil
}
func (o *ordersHandler) CreateOrder(c *gin.Context) {
	o.checkRequestFormat()
}

func (o *ordersHandler) GetOrders(c *gin.Context) {
	o.checkRequestFormat()
}

func NewOrdersHandler(processor ordersProcessor) *ordersHandler {
	return &ordersHandler{
		processor: processor,
	}
}
