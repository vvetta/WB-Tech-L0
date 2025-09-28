package usecase

import (
	"WB-Tech-L0/internal/domain"
	"context"
	"github.com/google/uuid"
)

type Cache interface {
	Set(key string, value *domain.Order) 
	Get(key string) (*domain.Order, bool)
}

type Repo interface {
	UpsertOrder(ctx context.Context, order *domain.Order) (bool, error)
	GetOrderById(ctx context.Context, orderId uuid.UUID) (*domain.Order, error)
	ListRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error)
}

type MessageConsumer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
