package repository

import (
	"WB-Tech-L0/internal/domain"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (p *PostgresRepo) UpsertOrder(ctx context.Context, order *domain.Order) error {

}

func (p *PostgresRepo) GetOrderById(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {

}

func (p *PostgresRepo) ListRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error) {

}
