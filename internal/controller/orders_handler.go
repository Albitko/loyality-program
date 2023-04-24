package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Albitko/loyalty-program/internal/entities"
	"github.com/Albitko/loyalty-program/internal/utils"
)

type ordersProcessor interface {
	CheckOrderExist(ctx context.Context, order int, userFromRequest string) error
	RegisterOrder(ctx context.Context, order int, userID string) error
	GetUserOrder(ctx context.Context, userID string) ([]entities.OrderWithTime, error)
}

type ordersHandler struct {
	processor ordersProcessor
}

func (o *ordersHandler) CreateOrder(c *gin.Context) {
	orderNumber, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	order, err := strconv.Atoi(string(orderNumber))
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{Message: err.Error()})
		return
	}
	if !utils.LuhnValid(order) {
		c.JSON(http.StatusUnprocessableEntity, entities.ErrorResponse{Message: "Invalid orders number"})
		return
	}
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	err = o.processor.CheckOrderExist(c, order, fmt.Sprintf("%v", userID))
	switch err {
	case entities.ErrNoOrderForUser:
		err = o.processor.RegisterOrder(c, order, fmt.Sprintf("%v", userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, entities.ErrorResponse{Message: "Order added"})
		return
	case entities.ErrOrderAlreadyCreatedByThisUser:
		c.JSON(http.StatusOK, entities.ErrorResponse{Message: "Order already registered by this user"})
		return
	case entities.ErrOrderAlreadyCreatedByAnotherUser:
		c.JSON(http.StatusConflict, entities.ErrorResponse{Message: "Order already registered by another user"})
		return
	default:
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}

}

func (o *ordersHandler) GetOrders(c *gin.Context) {
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	orders, err := o.processor.GetUserOrder(c, fmt.Sprintf("%v", userID))
	if errors.Is(err, entities.ErrNoOrderForUser) {
		c.JSON(http.StatusNoContent, entities.ErrorResponse{Message: "No orders for this user"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
	return
}

func NewOrdersHandler(processor ordersProcessor) *ordersHandler {
	return &ordersHandler{
		processor: processor,
	}
}
