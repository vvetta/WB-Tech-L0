package repository

import (
	"WB-Tech-L0/internal/domain"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"	
	"fmt"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (p *PostgresRepo) UpsertOrder(ctx context.Context, order *domain.Order) (created bool, err error) {
	if order == nil {
		//TODO Пока ошибка будет простой fmt, надо изменить.	
		return false, fmt.Errorf("Нулевой указатель заказа!")
	}

	return false, p.db.Transaction(func(tx *gorm.DB) error {
		for i := range order.Items {
			order.Items[i].OrderUID = order.OrderUID
		}

		result := tx.Session(&gorm.Session{
			FullSaveAssociations: true,
		}).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "order_uid"}},
			DoNothing: true,
		}).Create(order)

		if result.Error != nil {
			return result.Error
		}

		created = result.RowsAffected > 0
		
		return nil
	})
}

func (p *PostgresRepo) GetOrderById(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {

}

func (p *PostgresRepo) ListRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error) {

}
