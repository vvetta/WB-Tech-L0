package usecase

import (
	"WB-Tech-L0/internal/domain"
	"context"
	"time"
)

type OrderService struct {
	repo  Repo
	cache Cache
	log   Logger
}

func NewOrderService(r Repo, c Cache, l Logger) *OrderService {
	return &OrderService{repo: r, cache: c, log: l}
}

func (s *OrderService) GetByID(ctx context.Context, orderUID string) (*domain.Order, error) {
	if orderUID == "" {
		s.log.Error("usecase.GetByID: empty orderUID")
		return nil, domain.ErrNotFound
	}

	t0 := time.Now()
	if order, ok := s.cache.Get(orderUID); ok {
		t_cache := time.Since(t0)
		cacheMs := float64(t_cache) / float64(time.Millisecond)
		s.log.Debug("usecase.GetByID: cache hit", "order_uid", orderUID, "timeDuration_ms", cacheMs)
		return order, nil
	}
	
	tDB := time.Now()	
	order, err := s.repo.GetOrderById(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	t := time.Since(tDB)
	dbMs := float64(t) / float64(time.Millisecond)
	s.cache.Set(orderUID, order)
	s.log.Debug("usecase.GetByID: cache miss -> db", "order_uid", orderUID, "timeDuration_ms", dbMs)

	return order, nil
}

func (s *OrderService) WarmUpCache(ctx context.Context, limit int) error {
	s.log.Info("usecase.WarmUpCache: begin", "limit", limit)

	orders, err := s.repo.ListRecentOrders(ctx, limit)
	if err != nil {
		s.log.Error("usecase.WarmUpCache: repo error", "err", err)
		return err
	}

	for _, order := range orders {
		s.cache.Set(order.OrderUID, order)
	}

	s.log.Info("usecase.WarmUpCache: done", "count", len(orders))
	return nil
}
