package usecase

import (
	"WB-Tech-L0/domain"
	"context"
	"github.com/google/uuid"
)

type Cache interface {
	Set(key string, value *domain.Order) 
	Get(key string) (*domain.Order, bool)
}

type Repo interface {
	UpsertOrder(ctx context.Context, order *domain.Order) error
	GetOrderById(ctx context.Context, orderId uuid.UUID) (*domain.Order, error)
	ListRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error)
}
