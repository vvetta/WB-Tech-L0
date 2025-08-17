/*
Package models содержит в себе определение всех моделей приложения.
*/
package models

import (
	"time"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint `gorm:"primarykey" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Order struct {
	BaseModel
	OrderUID          string       `json:"order_uid" gorm:"primaryKey;unique"`
	TrackNumber       string       `json:"track_number"`
	Entry             string       `json:"entry"`
	Delivery          DeliveryInfo `json:"delivery" gorm:"embedded"`
	Payment           PaymentInfo  `json:"payment" gorm:"embedded"`
	Items             []ItemInfo   `json:"items" gorm:"foreignKey:OrderUID;references:OrderUID"`
	Locale            string       `json:"locale"`
	InternalSignature string       `json:"internal_signature"`
	CustomerID        string       `json:"customer_id"`
	DeliveryService   string       `json:"delivery_service"`
	ShardKey          string       `json:"shard_key"`
	SMID              int          `json:"sm_id"`
	DateCreated       string       `json:"date_created"`
	OOFShard          string       `json:"oof_shard"`
}

/*
Структуры DeliveryInfo и PaymentInfo просто встраиваем внутрь таблицы Order, а потом разворачиваем.
*/

type DeliveryInfo struct {
	Name    string `json:"name" gorm:"not null"`
	Phone   string `json:"phone" gorm:"not null"`
	Zip     string `json:"zip" gorm:"not null"`
	City    string `json:"city" gorm:"not null"`
	Address string `json:"address" gorm:"not null"`
	Region  string `json:"region" gorm:"not null"`
	Email   string `json:"email" gorm:"not null"`
}

type PaymentInfo struct {
	Transaction  string `json:"transaction" gorm:"not null"`
	RequestID    string `json:"request_id" gorm:"not null"`
	Currency     string `json:"currency" gorm:"not null"`
	Provider     string `json:"provider" gorm:"not null"`
	Amount       int    `json:"amount" gorm:"not null"`
	PaymentDT    int    `json:"payment_dt" gorm:"not null"`
	Bank         string `json:"bank" gorm:"not null"`
	DeliveryCost int    `json:"delivery_cost" gorm:"not null"`
	GoodsTotal   int    `json:"goods_total" gorm:"not null"`
	CustomFee    int    `json:"custom_fee" gorm:"not null"`
}

/*
Для товаров внутри заказа создаем отдельную таблицу, а потом связываем их с заказом по ключу заказа.
*/

type ItemInfo struct {
	BaseModel
	OrderUID string `json:"-" gorm:"index"`
	ChrtID      int    `json:"chrt_id" gorm:"not null"`
	TrackNumber string `json:"track_number" gorm:"not null"`
	Price       int    `json:"price" gorm:"not null"`
	RID         string `json:"rid" gorm:"not null"`
	Name        string `json:"name" gorm:"not null"`
	Sale        int    `json:"sale" gorm:"not null"`
	Size        string `json:"size" gorm:"not null"`
	TotalPrice  int    `json:"total_price" gorm:"not null"`
	NmID        int    `json:"nm_id" gorm:"not null"`
	Brand       string `json:"brand" gorm:"not null"`
	Status      int    `json:"status" gorm:"not null"`
}
