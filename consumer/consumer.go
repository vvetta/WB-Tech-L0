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
		log.Fatalf("Kafka brokers error: %v", err)
	}
	kafkaTopicName := config.GetKafkaTopicName()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kafkaBrokers,
		Topic:          kafkaTopicName,
		GroupID:        "wb-tech-l0-orders", // важно для идемпотентного чтения
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,                   // коммитим вручную после успеха
		// StartOffset:  kafka.LastOffset,   // <- расскомментируй в dev, чтобы пропускать старую историю при НОВОЙ группе
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Kafka reader close error: %v", err)
		}
	}()

	ctx := context.Background()

	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Kafka fetch error: %v", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Bad message (JSON): %v", err)
			if cerr := reader.CommitMessages(ctx, m); cerr != nil {
				log.Printf("Kafka commit error after bad msg: %v", cerr)
			}
			continue
		}

		created, err := database.AddOrderUpsert(&order)
		if err != nil {
			log.Printf("DB error: %v", err)
			continue
		}

		if created {
			log.Printf("Order %s inserted", order.OrderUID)
		} else {
			log.Printf("Order %s already exists — skipped", order.OrderUID)
		}

		if err := reader.CommitMessages(ctx, m); err != nil {
			log.Printf("Kafka commit error: %v", err)
		}
	}
}
