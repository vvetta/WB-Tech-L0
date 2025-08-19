package consumer

import (
	"log"
	"time"
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
		log.Fatalf("Ошибка при получении Brokers Kafka: %v", err)
	}
	kafkaTopicName := config.GetKafkaTopicName()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kafkaBrokers,
		Topic:          kafkaTopicName,
		GroupID:        "wb-tech-l0-orders",
		MinBytes:       1, // Благодаря этому параметру можно "Копить" сообщения.
		MaxBytes:       10e6,
		CommitInterval: 0, // Этот параметр отключает автоматические коммиты.
	})

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Kafka reader закрыт с ошибкой: %v", err)
		}
	}()

	ctx := context.Background()

	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Ошибка получения нового сообщения: %v", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Получен не валидный json заказа: %v", err)
			if cerr := reader.CommitMessages(ctx, m); cerr != nil {
				// Сообщение в любом случае коммитим, чтобы не читать его повторно.
				log.Printf("Ошибка коммита после получения не валидного json: %v", cerr)
			}
			continue
		}

		created, err := database.AddOrderUpsert(&order)
		if err != nil {
			log.Printf("Ошибка базы данных: %v", err)
			continue
		}

		if created {
			log.Printf("Заказ %s вставлен в базу данных.", order.OrderUID)
		} else {
			log.Printf("Заказ %s уже присутствует в базе данных — пропускаю", order.OrderUID)
		}

		if err := reader.CommitMessages(ctx, m); err != nil {
			log.Printf("Ошибка коммита: %v", err)
		}
	}
}
