package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Albitko/loyalty-program/internal/entities"
	"github.com/Albitko/loyalty-program/internal/utils"
)

type balanceProcessor interface {
	GetUserBalance(ctx context.Context, userID string) (entities.Balance, error)
	GetUserWithdrawals(ctx context.Context, userID string) ([]entities.WithdrawWithTime, error)
	Withdraw(ctx context.Context, userID string, request entities.Withdraw) error
}
type balanceHandler struct {
	processor balanceProcessor
}

func (b *balanceHandler) GetBalance(c *gin.Context) {
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	balance, err := b.processor.GetUserBalance(c, fmt.Sprintf("%v", userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (b *balanceHandler) Withdraw(c *gin.Context) {
	var request entities.Withdraw
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	orderNumber, err := strconv.Atoi(request.Order)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{Message: err.Error()})
		return
	}
	if !utils.LuhnValid(orderNumber) {
		c.JSON(http.StatusUnprocessableEntity, entities.ErrorResponse{Message: "Wrong order number"})
		return
	}

	err = b.processor.Withdraw(c, fmt.Sprintf("%v", userID), request)
	if errors.Is(err, entities.ErrInsufficientFunds) {
		c.JSON(http.StatusPaymentRequired, entities.ErrorResponse{Message: "Insufficient funds"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
}

func (b *balanceHandler) GetWithdrawn(c *gin.Context) {
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	withdrawals, err := b.processor.GetUserWithdrawals(c, fmt.Sprintf("%v", userID))
	if errors.Is(err, entities.ErrNoWithdrawals) {
		c.JSON(http.StatusNoContent, entities.ErrorResponse{Message: err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, withdrawals)
}

func NewBalanceHandler(processor balanceProcessor) *balanceHandler {
	return &balanceHandler{
		processor: processor,
	}
}
