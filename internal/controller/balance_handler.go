package controller

import (
	"context"

	"github.com/gin-gonic/gin"
)

type balanceProcessor interface {
	GetUserBalance(ctx context.Context) error
	GetUserWithdrawals(ctx context.Context) error
	CheckUserAvailableBalance(ctx context.Context) (int, error)
	Withdraw(ctx context.Context) error
}
type balanceHandler struct {
	processor balanceProcessor
}

func (b *balanceHandler) GetBalance(c *gin.Context) {
}

func (b *balanceHandler) Withdraw(c *gin.Context) {
}

func (b *balanceHandler) GetWithdrawn(c *gin.Context) {
}

func NewBalanceHandler(processor balanceProcessor) *balanceHandler {
	return &balanceHandler{
		processor: processor,
	}
}
