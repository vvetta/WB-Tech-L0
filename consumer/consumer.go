package consumer

import (
	"log"
	"fmt"
	"context"
	"encoding/json"

	"WB-Tech-L0/config"
	"WB-Tech-L0/models"
	"WB-Tech-L0/database"

	"github.com/segmentio/kafka-go"
)

func ConsumeOrders() {
	kafkaBrokers, err := config.GetKafkaBrokers()
	if err != nil {
		fmt.Print(err)
	}

	kafkaTopicName := config.GetKafkaTopicName()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic: kafkaTopicName,
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
			log.Println("Получено не валидное сообщение: ", err)
			continue
		}
	
		//TODO Проверка наличия сообщения в базе данных.

		database.AddOrderToDB(order)

		fmt.Println("From kafka: ", order.OrderUID)
	}
}
