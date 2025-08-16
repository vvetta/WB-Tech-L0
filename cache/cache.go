package cache

import (
	"sync"

	"WB-Tech-L0/models"
)

var (
	orderCache sync.Map
)

func Set(order *models.Order) {
// Добавляет заказ в map.
	if order == nil || order.OrderUID == "" {
		return
	}
	orderCache.Store(order.OrderUID, order)
}

func Get(order_uid string) *models.Order {
// Получаем заказ из кеша.
	value, ok := orderCache.Load(order_uid)
	if !ok {
		return nil
	}
	return value.(*models.Order)
}

