package usecase

import (
	"context"
)

type balanceProcessor struct {
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

func NewBalanceProcessor() *balanceProcessor {
	return &balanceProcessor{}
}
