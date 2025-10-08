package usecase

import (
	"WB-Tech-L0/internal/domain"
	"context"
)

type OrderReader interface {
	GetByID(ctx context.Context, orderUID string) (*domain.Order, error)
}

type CacheWarmer interface {
	WarmUpCache(ctx context.Context, limit int) error
}

type Cache interface {
	Set(key string, value *domain.Order)
	Get(key string) (*domain.Order, bool)
}

type Repo interface {
	UpsertOrder(ctx context.Context, order *domain.Order) (bool, error)
	GetOrderById(ctx context.Context, orderUID string) (*domain.Order, error)
	ListRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error)
}

type MessageConsumer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
