package consumer

import (
	"fmt"
	"WB-Tech-L0/models"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

func ConsumeOrders() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "test-topic",
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var order models.Order
		err = json.Unmarshal(msg.Value, &order)
		if err != nil {
			continue
		}

		// Тут будет сохранение в базу данных, на данный момент просто Print.
		fmt.Println(order)
	}
}
