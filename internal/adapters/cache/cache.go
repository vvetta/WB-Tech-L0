package cache

import (
	"sync"
	"WB-Tech-L0/internal/domain"
)

type MemoryCache struct {
	mu sync.RWMutex
	store map[string]*domain.Order
	itemLimit int 
	deleteItemCount int
}

func NewMemoryCache(limit, deleteCount int) *MemoryCache {
	return &MemoryCache{
		store: make(map[string]*domain.Order),
		itemLimit: limit,
		deleteItemCount: deleteCount,
	}
}

func (m *MemoryCache) Set(key string, value *domain.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Псевдо защита от переполнения памяти.
	// Конечно это далеко не продакшн. :)
	if len(m.store) >= m.itemLimit {
		var keys []string

		keysCount := 0

		for key, _ := range m.store {
			// Получаем N ключей рандомных элементов.
			if keysCount == m.deleteItemCount { break }	

			keys = append(keys, key)
			keysCount++
		}

		for _, key := range keys {
			// Проходимся по ключам и удаляем.
			delete(m.store, key)
		}
	}

	m.store[key] = value
}

func (m *MemoryCache) Get(key string) (*domain.Order, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	order, ok := m.store[key]
	return order, ok
}
