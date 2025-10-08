package repository

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Order struct {
	BaseModel
	OrderUID          string `gorm:"primaryKey;unique"`
	TrackNumber       string
	Entry             string
	Delivery          DeliveryInfo `gorm:"embedded"`
	Items             []ItemInfo   `gorm:"foreignKey:OrderUID;references:OrderUID"`
	Payment           PaymentInfo  `gorm:"embedded"`
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	DeliveryInfo      string
	ShardKey          string
	SMID              int
	DateCreated       string
	OOFShard          string
}

type DeliveryInfo struct {
	Name    string `gorm:"not null"`
	Phone   string `gorm:"not null"`
	Zip     string `gorm:"not null"`
	City    string `gorm:"not null"`
	Address string `gorm:"not null"`
	Region  string `gorm:"not null"`
	Email   string `gorm:"not null"`
}

type PaymentInfo struct {
	Transaction  string `gorm:"not null"`
	RequestID    string `gorm:"not null"`
	Currency     string `gorm:"not null"`
	Provider     string `gorm:"not null"`
	Amount       int    `gorm:"not null"`
	PaymentDT    int    `gorm:"not null"`
	Bank         string `gorm:"not null"`
	DeliveryCost int    `gorm:"not null"`
	GoodsTotal   int    `gorm:"not null"`
	CustomFee    int    `gorm:"not null"`
}

type ItemInfo struct {
	BaseModel
	OrderUID    string `gorm:"index"`
	ChrtID      int    `gorm:"not null"`
	TrackNumber string `gorm:"not null"`
	Price       int    `gorm:"not null"`
	RID         string `gorm:"not null"`
	Name        string `gorm:"not null"`
	Sale        int    `gorm:"not null"`
	Size        string `gorm:"not null"`
	TotalPrice  int    `gorm:"not null"`
	NmID        int    `gorm:"not null"`
	Brand       string `gorm:"not null"`
	Status      int    `gorm:"not null"`
}
