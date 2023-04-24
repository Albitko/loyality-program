package usecase

import (
	"context"

	"github.com/Albitko/loyalty-program/internal/entities"
)

type balanceRepository interface {
	GetUserBalance(ctx context.Context, user string) (float64, error)
	GetUserWithdrawn(ctx context.Context, user string) (float64, error)
	GetUserAllWithdrawals(ctx context.Context, userId string) ([]entities.WithdrawWithTime, error)
	Withdraw(ctx context.Context, userID string, withdrawRequest entities.Withdraw) error
}

type balanceProcessor struct {
	repository balanceRepository
}

func (b balanceProcessor) GetUserBalance(ctx context.Context, userId string) (entities.Balance, error) {
	var balance entities.Balance

	accrualsTotal, err := b.repository.GetUserBalance(ctx, userId)
	if err != nil {
		return balance, err
	}
	withdrawnTotal, err := b.repository.GetUserWithdrawn(ctx, userId)
	if err != nil {
		return balance, err
	}
	balance.Current = accrualsTotal - withdrawnTotal
	balance.Withdrawn = withdrawnTotal

	return balance, nil
}

func (b balanceProcessor) GetUserWithdrawals(ctx context.Context, userId string) ([]entities.WithdrawWithTime, error) {
	withdrawals, err := b.repository.GetUserAllWithdrawals(ctx, userId)
	return withdrawals, err
}

func (b balanceProcessor) Withdraw(ctx context.Context, userID string, request entities.Withdraw) error {
	balance, err := b.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance.Current < request.Sum {
		return entities.ErrInsufficientFunds
	}
	err = b.repository.Withdraw(ctx, userID, request)
	if err != nil {
		return err
	}
	return nil
}

func NewBalanceProcessor(repository balanceRepository) *balanceProcessor {
	return &balanceProcessor{
		repository: repository,
	}
}
