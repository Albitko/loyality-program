package usecase

import (
	"context"

	"github.com/Albitko/loyalty-program/internal/entities"
)

type ordersRepository interface {
	GetUserForOrder(ctx context.Context, order string) (string, error)
	CreateOrder(ctx context.Context, order string) error
	GetOrdersForUser(ctx context.Context, user string) ([]string, error)
}

type ordersQueue interface {
	Push(entities.Order)
}

type ordersProcessor struct {
	repository ordersRepository
	queue      ordersQueue
}

func (o *ordersProcessor) CheckOrderExist(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordersProcessor) RegisterOrder(ctx context.Context) error {
	//TODO implement me
	// write order to DB with CreateOrder
	// push it to the queue
	panic("implement me")
}

func (o *ordersProcessor) GetUserOrder(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewOrdersProcessor(repository ordersRepository, queue ordersQueue) *ordersProcessor {
	return &ordersProcessor{
		repository: repository,
		queue:      queue,
	}
}
