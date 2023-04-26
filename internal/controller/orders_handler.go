package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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
		utils.Logger.Error("ordersHandler:CreateOrder - read body request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	order, err := strconv.Atoi(string(orderNumber))
	if err != nil {
		utils.Logger.Error("ordersHandler:CreateOrder - convert order to int", zap.Error(err))
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{Message: err.Error()})
		return
	}
	if !utils.LuhnValid(order) {
		utils.Logger.Error("ordersHandler:CreateOrder - Luhn validation", zap.Int("order", order))
		c.JSON(http.StatusUnprocessableEntity, entities.ErrorResponse{Message: "Invalid orders number"})
		return
	}
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		utils.Logger.Error("ordersHandler:CreateOrder - extract userID", zap.Bool("isExtract", isExtract))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	err = o.processor.CheckOrderExist(c, order, fmt.Sprintf("%v", userID))
	switch err {
	case entities.ErrNoOrderForUser:
		err = o.processor.RegisterOrder(c, order, fmt.Sprintf("%v", userID))
		if err != nil {
			utils.Logger.Error("ordersHandler:CreateOrder - register order", zap.Error(err))
			c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, entities.ErrorResponse{Message: "Order added"})
		return
	case entities.ErrOrderAlreadyCreatedByThisUser:
		utils.Logger.Error("ordersHandler:CreateOrder - order already created", zap.Error(err))
		c.JSON(http.StatusOK, entities.ErrorResponse{Message: "Order already registered by this user"})
		return
	case entities.ErrOrderAlreadyCreatedByAnotherUser:
		utils.Logger.Error("ordersHandler:CreateOrder - order already created", zap.Error(err))
		c.JSON(http.StatusConflict, entities.ErrorResponse{Message: "Order already registered by another user"})
		return
	default:
		utils.Logger.Error("ordersHandler:CreateOrder - check order exist", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}

}

func (o *ordersHandler) GetOrders(c *gin.Context) {
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		utils.Logger.Error("ordersHandler:GetOrders - extract userID", zap.Bool("isExtract", isExtract))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	orders, err := o.processor.GetUserOrder(c, fmt.Sprintf("%v", userID))
	if errors.Is(err, entities.ErrNoOrderForUser) {
		utils.Logger.Error("ordersHandler:GetOrders - no orders for user", zap.Error(err))
		c.JSON(http.StatusNoContent, entities.ErrorResponse{Message: "No orders for this user"})
		return
	}
	if err != nil {
		utils.Logger.Error("ordersHandler:GetOrders - GetUserOrder", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func NewOrdersHandler(processor ordersProcessor) *ordersHandler {
	return &ordersHandler{
		processor: processor,
	}
}
