package repository

import (
	"WB-Tech-L0/internal/domain"
	"context"
	"fmt"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (p *PostgresRepo) UpsertOrder(ctx context.Context, order *domain.Order) (bool, error) {
	if order == nil {
		return false, fmt.Errorf("Нулевой указатель заказа!")
	}
	
	var created bool
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gm := toGORMOrder(order)

		res := tx.Session(&gorm.Session{FullSaveAssociations: true}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "order_uid"}},
			DoNothing: true,
		}).
		Create(gm)
		
		if res.Error != nil {
			return res.Error
		}

		created = res.RowsAffected > 0 
		return nil
	})

	return created, err
}

func (p *PostgresRepo) GetOrderById(ctx context.Context, orderUID string) (*domain.Order, error) {
	if orderUID == "" {
		return nil, fmt.Errorf("Передан пустой id заказа!")
	}

	var gormOrder Order

	err := p.db.Preload("Items").Where("order_uid = ?", orderUID).First(&gormOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Заказ с данным ( %s )id не найден!", orderUID)	
		}
		return nil, fmt.Errorf("Произошла ошибка при получении заказа! %v", err)
	}

	domainOrder := toDomainOrder(&gormOrder)

	return domainOrder, nil
}

func (p *PostgresRepo) ListRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error) {
	var gormOrders []Order

	err := p.db.Preload("Delivery").Preload("Payment").Preload("Items").Order("date_created DESC").Limit(limit).Find(&gormOrders).Error
	if err != nil {
		return nil, fmt.Errorf("Произошла ошибка при получении списка заказов! %v", err)
	}

	domainOrders := make([]*domain.Order, len(gormOrders))
	for i, gormOrder := range gormOrders {
		domainOrders[i] = toDomainOrder(&gormOrder)
	}

	return domainOrders, nil
}

