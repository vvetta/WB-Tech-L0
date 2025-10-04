package cache

import (
	"sync"

	"WB-Tech-L0/internal/domain"
	"WB-Tech-L0/internal/usecase"
)

type MemoryCache struct {
	mu sync.RWMutex
	store map[string]*domain.Order
	itemLimit int 
	deleteItemCount int
	log usecase.Logger
}

func NewMemoryCache(limit, deleteCount int, log usecase.Logger) *MemoryCache {
	return &MemoryCache{
		store: make(map[string]*domain.Order),
		itemLimit: limit,
		deleteItemCount: deleteCount,
		log: log,
	}
}

func (m *MemoryCache) Set(key string, value *domain.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Debug("cache.Set: begin", "key", key)
	
	// Псевдо защита от переполнения памяти.
	// Конечно это далеко не продакшн. :)
	if len(m.store) >= m.itemLimit {

		m.log.Debug("cache.Set: overflow protection activated", "itemLimit", m.itemLimit)

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

	m.log.Info("cache.Set: created", "key", key)

	m.store[key] = value
}

func (m *MemoryCache) Get(key string) (*domain.Order, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.log.Debug("cache.Get")

	order, ok := m.store[key]
	return order, ok
}
