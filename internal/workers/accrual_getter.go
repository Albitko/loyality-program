package workers

import (
	"github.com/Albitko/loyalty-program/internal/entities"
)

type accrualGetter struct {
	accrualURL string
}

func (s *accrualGetter) GetAccrual(orderID string) (entities.Order, error) {
	var order entities.Order
	return order, nil
}

func newAccrualGetter(accrualURL string) *accrualGetter {
	return &accrualGetter{accrualURL: accrualURL}
}
