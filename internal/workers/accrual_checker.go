package workers

import (
	"context"
	"runtime"

	"github.com/Albitko/loyalty-program/internal/entities"
	"github.com/Albitko/loyalty-program/internal/repo"
)

type ordersQueue interface {
	PopWait() entities.Order
	Push(entities.Order)
}

type orderStorage interface {
	UpdateOrder(context.Context, entities.Order) error
}

type accrualChecker struct {
	queue   ordersQueue
	storage orderStorage
	getter  *accrualGetter
	ctx     context.Context
}

func (a *accrualChecker) loop() {
	for {
		order := a.queue.PopWait()
		updatedOrder, err := a.getter.GetAccrual(order.OrderID)
		if err != nil {
			continue
		}
		err = a.storage.UpdateOrder(a.ctx, updatedOrder)
		if err != nil {
			continue
		}
		if updatedOrder.Status != "INVALID" && updatedOrder.Status != "PROCESSED" {
			a.queue.Push(updatedOrder)
		}
	}
}
func newAccrualChecker(
	ctx context.Context, storage orderStorage, queue ordersQueue, getter *accrualGetter,
) *accrualChecker {
	return &accrualChecker{
		ctx:     ctx,
		queue:   queue,
		storage: storage,
		getter:  getter,
	}
}

func InitWorkers(ctx context.Context, storage orderStorage, accrualServiceURL string) ordersQueue {
	queue := repo.NewQueue()
	checkers := make([]*accrualChecker, 0, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		checkers = append(checkers, newAccrualChecker(ctx, storage, queue, newAccrualGetter(accrualServiceURL)))
	}

	for _, checker := range checkers {
		go checker.loop()
	}
	return queue
}
