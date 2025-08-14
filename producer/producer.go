package producer

import (
	"fmt"
	"time"
	"context"
	"math/rand"
	"encoding/json"

	"WB-Tech-L0/models"
	"WB-Tech-L0/config"

	"github.com/segmentio/kafka-go"
)

func AddTestMessageToKafka() error {
	ctx := context.Background()

	kafkaBrokers, err := config.GetKafkaBrokers()
	if err != nil {
		return err
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: kafkaBrokers,
		Topic: "test-topic",
	})

	defer writer.Close()

	jsonOrder, jsonOrderId, err := generateRandomOrderJson()
	if err != nil {
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Value: jsonOrder,
	})

	if err != nil {
		return err
	}

	fmt.Println("OrderUID: ", jsonOrderId)

	return nil
}

func PrintTestJsonOrder() {
	jsonOrder, jsonOrderId, err := generateRandomOrderJson()
	if err != nil {
		fmt.Printf("Ошибка при выводе тестового заказа: %v", err)
	}

	order := models.Order{}
	json.Unmarshal(jsonOrder, &order)

	fmt.Println("OrderID: ", jsonOrderId)
	fmt.Println(order)
}

func generateRandomOrderJson() ([]byte, string, error) {
// Генерирует заказ с рандомными значениями для теста.

	items_count := getRandomInt()
	items := []models.ItemInfo{}

	for i := 0; i < items_count; i++ {
		item := models.ItemInfo{
			ChrtID: getRandomInt(),
			TrackNumber: generateRandomString(getRandomInt()),
			Price: getRandomInt(),
			RID: generateRandomString(getRandomInt()),
			Name: generateRandomString(getRandomInt()),
			Sale: getRandomInt(),
			Size: generateRandomString(getRandomInt()),
			TotalPrice: getRandomInt(),
			NmID: getRandomInt(),
			Brand: generateRandomString(getRandomInt()),
			Status: getRandomInt(),
		}

		items = append(items, item)
	}

	payment := models.PaymentInfo{
		Transaction: generateRandomString(getRandomInt()),
		RequestID: generateRandomString(getRandomInt()),
		Currency: generateRandomString(getRandomInt()),
		Provider: generateRandomString(getRandomInt()),
		Amount: getRandomInt(),
		PaymentDT: getRandomInt(),
		Bank: generateRandomString(getRandomInt()),
		DeliveryCost: getRandomInt(),
		GoodsTotal: getRandomInt(),
		CustomFee: getRandomInt(),
	}

	delivery := models.DeliveryInfo{
		Name: generateRandomString(getRandomInt()),
		Phone: generateRandomString(getRandomInt()),
		Zip: generateRandomString(getRandomInt()),
		City: generateRandomString(getRandomInt()),
		Address: generateRandomString(getRandomInt()),
		Region: generateRandomString(getRandomInt()),
		Email: generateRandomString(getRandomInt()),	
	}

	order := models.Order{
		OrderUID: generateRandomString(getRandomInt()),
		TrackNumber: generateRandomString(getRandomInt()),
		Entry: generateRandomString(getRandomInt()),
		Delivery: delivery,
		Payment: payment,
		Locale: generateRandomString(getRandomInt()),
		InternalSignature: generateRandomString(getRandomInt()),
		CustomerID: generateRandomString(getRandomInt()),
		DeliveryService: generateRandomString(getRandomInt()),
		ShardKey: generateRandomString(getRandomInt()),
		SMID: getRandomInt(),
		DateCreated: generateRandomString(getRandomInt()),
		OOFShard: generateRandomString(getRandomInt()),
		Items: items,
	}

	b, err := json.Marshal(order)

	if err != nil {
		return b, order.OrderUID, err
	}

	return b, order.OrderUID, nil
}

func getRandomInt() int {
// Возвращает рандомное значение.
	return rand.Intn(100)
}

func generateRandomString(string_length int) string {
	// Вспомогательная функция для генерирования строки.

	const charset = "zaqwsxcderfvbgtyhnmjuiklop"
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, string_length)

	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}

	return string(b)
}
