package workers

import (
	"github.com/Albitko/loyalty-program/internal/entities"
	"github.com/Albitko/loyalty-program/internal/utils"
)

type accrualGetter struct {
	accrualURL string
}

func (s *accrualGetter) GetAccrual(orderID string) (entities.Order, error) {
	var order entities.Order
	_, err := utils.RestyClient.R().
		EnableTrace().
		SetResult(&order).
		Get(s.accrualURL + "/api/orders/" + orderID)

	if err != nil {
		return order, err
	}
	return order, nil
}

func newAccrualGetter(accrualURL string) *accrualGetter {
	return &accrualGetter{accrualURL: accrualURL}
}
