package usecase

import (
	"context"
)

type balanceRepository interface {
	GetUserBalance(ctx context.Context, user string) (string, error)
	GetUserWithdrawn(ctx context.Context, user string) (string, error)
	Withdraw(ctx context.Context, amount string) error
}

type balanceProcessor struct {
	repository balanceRepository
}

func (b balanceProcessor) GetUserBalance(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b balanceProcessor) GetUserWithdrawals(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b balanceProcessor) CheckUserAvailableBalance(ctx context.Context) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (b balanceProcessor) Withdraw(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewBalanceProcessor(repository balanceRepository) *balanceProcessor {
	return &balanceProcessor{
		repository: repository,
	}
}
