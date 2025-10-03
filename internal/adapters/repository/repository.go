package repository

import (
	"WB-Tech-L0/internal/domain"
	"WB-Tech-L0/internal/usecase"
	
	"context"
	"fmt"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresRepo struct {
	db *gorm.DB
	log usecase.Logger
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (p *PostgresRepo) UpsertOrder(ctx context.Context, order *domain.Order) (bool, error) {
	if order == nil {
		err := fmt.Errorf("Нулевой указатель заказа!")
		if p.log != nil {
			p.log.Error("repo.upsertOrder: invalid input", "err", err)
		}
		return false, err
	}

	if p.log != nil {
		p.log.Debug("repo.upsertOrder: begin", "order_uid", order.OrderUID)
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

	if err != nil {
		if p.log != nil {
			p.log.Error("repo.upsertOrder: db error", "order_uid", order.OrderUID)
		}
		return false, fmt.Errorf("Ошибка сохранения заказа в базу данных: %w", err)
	}

	if p.log != nil {
		if created {
			p.log.Info("repo.upsertOrder: created", "order_uid", order.OrderUID)
		} else {
			p.log.Info("repo.upsertOrder: already exists", "order_uid", order.OrderUID)
		}
	}

	return created, err
}

func (p *PostgresRepo) GetOrderById(ctx context.Context, orderUID string) (*domain.Order, error) {
	if orderUID == "" {
		err := fmt.Errorf("Передан пустой id заказа!")
		if p.log != nil {
			p.log.Error("repo.GetOrderById: invalid input", "err", err)	
		}
		return nil, err
	}

	if p.log != nil {
		p.log.Debug("repo.GetOrderById: begin", "order_uid", orderUID)
	}

	var gormOrder Order
	err := p.db.Preload("Items").Where("order_uid = ?", orderUID).First(&gormOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if p.log != nil {
				p.log.Info("repo.GetOrderById: not found", "order_uid", orderUID)
			}
			return nil, fmt.Errorf("Заказ с данным ( %s )id не найден!", orderUID)	
		}
		if p.log != nil {
			p.log.Error("repo.GetOrderById: db error", "order_uid", orderUID)
		}
		return nil, fmt.Errorf("Произошла ошибка при получении заказа! %w", err)
	}

	domainOrder := toDomainOrder(&gormOrder)

	if p.log != nil {
		p.log.Debug("repo.GetOrderById: success", "order_uid", orderUID)
	}

	return domainOrder, nil
}

func (p *PostgresRepo) ListRecentOrders(ctx context.Context, limit int) ([]*domain.Order, error) {
	if p.log != nil {
		p.log.Debug("repo.ListRecentOrders: begin", "limit", limit)
	}

	var gormOrders []Order

	err := p.db.Preload("Delivery").Preload("Payment").Preload("Items").Order("date_created DESC").Limit(limit).Find(&gormOrders).Error
	if err != nil {
		if p.log != nil {
			p.log.Error("repo.ListRecentOrders: db error", "err", err)
		}
		return nil, fmt.Errorf("Произошла ошибка при получении списка заказов! %v", err)
	}

	domainOrders := make([]*domain.Order, len(gormOrders))
	for i, gormOrder := range gormOrders {
		domainOrders[i] = toDomainOrder(&gormOrder)
	}

	if p.log != nil {
		p.log.Info("repo.ListRecentOrders: success", "count", len(domainOrders))
	}

	return domainOrders, nil
}

