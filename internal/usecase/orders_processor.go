package usecase

import (
	"context"
	"strconv"

	"github.com/Albitko/loyalty-program/internal/entities"
)

//go:generate mockery --name ordersRepository
type ordersRepository interface {
	GetUserForOrder(ctx context.Context, order string) (string, error)
	CreateOrder(ctx context.Context, order entities.Order, userID string) error
	GetOrdersForUser(ctx context.Context, user string) ([]entities.OrderWithTime, error)
}

//go:generate mockery --name ordersQueue
type ordersQueue interface {
	Push(entities.Order)
}

type ordersProcessor struct {
	repository ordersRepository
	queue      ordersQueue
}

func (o *ordersProcessor) CheckOrderExist(ctx context.Context, order int, userFromRequest string) error {
	userFromDB, err := o.repository.GetUserForOrder(ctx, strconv.Itoa(order))
	if err != nil {
		return err
	}

	if userFromDB == userFromRequest {
		return entities.ErrOrderAlreadyCreatedByThisUser
	}
	return entities.ErrOrderAlreadyCreatedByAnotherUser
}

func (o *ordersProcessor) RegisterOrder(ctx context.Context, orderNumber int, userID string) error {
	var order entities.Order

	order.OrderID = strconv.Itoa(orderNumber)
	order.Status = "NEW"

	err := o.repository.CreateOrder(ctx, order, userID)
	if err != nil {
		return err
	}
	o.queue.Push(order)
	return nil
}

func (o *ordersProcessor) GetUserOrder(ctx context.Context, userID string) ([]entities.OrderWithTime, error) {
	orders, err := o.repository.GetOrdersForUser(ctx, userID)
	if err != nil {
		return orders, err
	}

	return orders, nil
}

func NewOrdersProcessor(repository ordersRepository, queue ordersQueue) *ordersProcessor {
	return &ordersProcessor{
		repository: repository,
		queue:      queue,
	}
}
