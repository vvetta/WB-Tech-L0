package database


import (
	
	"WB-Tech-L0/models"
  "WB-Tech-L0/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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


func AddOrderToDB(order models.Order) error {
// Сохраняет заказ в базу данных.

	return nil
}


func GetOrdersFromDB(orders_uid []string) ([]models.Order, error) {
// Возвращает заказы из базы данных.
	var orders []models.Order

	return orders, nil
}

