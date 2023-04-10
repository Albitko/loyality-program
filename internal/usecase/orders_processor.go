package usecase

import (
	"context"
)

type ordersProcessor struct {
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

func NewOrdersProcessor() *ordersProcessor {
	return &ordersProcessor{}
}
