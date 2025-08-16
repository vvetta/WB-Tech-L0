package cache

import (
	"sync"

	"WB-Tech-L0/models"
)

var (
	orderCache sync.Map
)

func Set(order *models.Order) {
	orderCache.Store(order.OrderUID, order)
}

func Get(order_uid string) *models.Order {
	value, ok := orderCache.Load(order_uid)
	if !ok {
		return nil
	}
	return value.(*models.Order)
}
