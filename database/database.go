package database


import (
	"fmt"	

	"WB-Tech-L0/models"
  "WB-Tech-L0/config"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func InitDB() error {
// Инициализация базы данных.
	dsn := config.GetDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	//TODO Нужно сделать автомиграцию.
	err = DB.AutoMigrate(&models.Order{}, &models.ItemInfo{})
	if err != nil {
		return err
	}

	return nil
}


func AddOrderUpsert(order *models.Order) (created bool, err error) {
    if order == nil {
        return false, fmt.Errorf("Нулевой указатель заказа")
    }

    return created, DB.Transaction(func(tx *gorm.DB) error {
        // убедимся, что дочкам проставлен FK
        for i := range order.Items {
            order.Items[i].OrderUID = order.OrderUID
        }

        res := tx.Session(&gorm.Session{FullSaveAssociations: true}).
            Clauses(clause.OnConflict{
                Columns:   []clause.Column{{Name: "order_uid"}},
                DoNothing: true, // или DoUpdates при нужде обновлять поля
            }).
            Create(order)

        if res.Error != nil {
            return res.Error
        }

        created = res.RowsAffected > 0
        return nil
    })
}


func GetOrdersFromDB(orders_uid []string, opts ...string) ([]models.Order, error) {
// Возвращает заказы из базы данных.
	if len(opts) != 0 && opts[0] == "all" {
		/*
		В этом случае выводятся все записи из базы данных, 
		это нужно для первичного заполнения кеша.
		*/
		
		var allOrders []models.Order
		err := DB.Preload("Items").Find(&allOrders).Error
		if err != nil {
			return nil, fmt.Errorf("Ошибка при получении всех заказов! %v", err)
		}

		return allOrders, nil
	} else {

		if len(orders_uid) == 0 {
			return nil, fmt.Errorf("Пустой список id заказов!")
		}

		var orders []models.Order

		err := DB.Preload("Items").Where("order_uid IN ?", orders_uid).Find(&orders).Error
		if err != nil {
			return nil, fmt.Errorf("Ошибка при получении заказов: %v", err)
		}

		if len(orders) == 0 {
			return nil, fmt.Errorf("Заказы не найдены!")
		}

		return orders, nil
	}
}

