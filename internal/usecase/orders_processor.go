package usecase

import (
	"context"
)

type ordersRepository interface {
	GetUserForOrder(ctx context.Context, order string) (string, error)
	CreateOrder(ctx context.Context, order string) error
	GetOrdersForUser(ctx context.Context, user string) ([]string, error)
}

type ordersProcessor struct {
	repository ordersRepository
}

func (o *ordersProcessor) CheckOrderExist(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordersProcessor) RegisterOrder(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordersProcessor) GetUserOrder(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewOrdersProcessor(repository ordersRepository) *ordersProcessor {
	return &ordersProcessor{
		repository: repository,
	}
}
