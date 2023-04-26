package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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
		utils.Logger.Error("balanceHandler:GetBalance extract userID", zap.Bool("isExtract", isExtract))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	balance, err := b.processor.GetUserBalance(c, fmt.Sprintf("%v", userID))
	if err != nil {
		utils.Logger.Error("balanceHandler:GetBalance from balanceProcessor", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, balance)
}

func (b *balanceHandler) Withdraw(c *gin.Context) {
	var request entities.Withdraw
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		utils.Logger.Error("balanceHandler:Withdraw - extract userID", zap.Bool("isExtract", isExtract))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		utils.Logger.Error("balanceHandler:Withdraw - request bind JSON ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
	orderNumber, err := strconv.Atoi(request.Order)
	if err != nil {
		utils.Logger.Error("balanceHandler:Withdraw - convert order from string to int", zap.Error(err))
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{Message: err.Error()})
		return
	}
	if !utils.LuhnValid(orderNumber) {
		utils.Logger.Error(
			"balanceHandler:Withdraw - Luhn check failed", zap.Int("orderNumber", orderNumber),
		)
		c.JSON(http.StatusUnprocessableEntity, entities.ErrorResponse{Message: "Wrong order number"})
		return
	}

	err = b.processor.Withdraw(c, fmt.Sprintf("%v", userID), request)
	if errors.Is(err, entities.ErrInsufficientFunds) {
		utils.Logger.Error("balanceHandler:Withdraw - insufficient funds", zap.Error(err))
		c.JSON(http.StatusPaymentRequired, entities.ErrorResponse{Message: "Insufficient funds"})
		return
	}
	if err != nil {
		utils.Logger.Error("balanceHandler:Withdraw - balanceProcessor error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: err.Error()})
		return
	}
}

func (b *balanceHandler) GetWithdrawn(c *gin.Context) {
	userID, isExtract := c.Get("x-user-id")
	if !isExtract {
		utils.Logger.Error("balanceHandler:GetWithdrawn - extract userID", zap.Bool("isExtract", isExtract))
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Message: "Invalid x-user-id"})
		return
	}
	withdrawals, err := b.processor.GetUserWithdrawals(c, fmt.Sprintf("%v", userID))
	if errors.Is(err, entities.ErrNoWithdrawals) {
		utils.Logger.Error("balanceHandler:GetWithdrawn - no withdrawals for user", zap.Error(err))
		c.JSON(http.StatusNoContent, entities.ErrorResponse{Message: err.Error()})
		return
	}
	if err != nil {
		utils.Logger.Error("balanceHandler:GetWithdrawn - GetUserWithdrawals err", zap.Error(err))
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
