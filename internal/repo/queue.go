package repo

import (
	"github.com/Albitko/loyalty-program/internal/entities"
)

type queue struct {
	ch chan entities.Order
}

func (q *queue) Push(t entities.Order) {
	q.ch <- t
}

func (q *queue) PopWait() entities.Order {
	return <-q.ch
}

func NewQueue() *queue {
	return &queue{
		ch: make(chan entities.Order, 15),
	}
}
