package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"

	"WB-Tech-L0/internal/usecase"
)

type Config struct {
	Brokers        []string
	Topic          string
	GroupID        string
	MinBytes       int
	MaxBytes       int
	StartOffset    int64
	CommitInterval time.Duration
}

type Consumer struct {
	reader *kafkago.Reader
	log    usecase.Logger
	repo   usecase.Repo
	cache  usecase.Cache
}

func New(cfg Config, repo usecase.Repo, cache usecase.Cache, log usecase.Logger) (*Consumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("")
	}

	if cfg.Topic == "" || cfg.GroupID == "" {
		return nil, fmt.Errorf("")
	}

	rd := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		StartOffset:    kafkago.FirstOffset,
		CommitInterval: cfg.CommitInterval,
	})

	return &Consumer{
		reader: rd,
		repo:   repo,
		cache:  cache,
		log:    log,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {

	c.log.Info("kafkaConsumer.Start: begin")
	defer c.log.Info("kafkaConsumer.Start: stop")

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			c.log.Error("kafkaConsumer.Start: fetch error", "err", err)
			//TODO Тут можно поставить задержку. Пока не буду ставить.
			continue
		}

		if err := c.process(ctx, m.Value); err != nil {
			c.log.Error("kafkaConsumer.Start: handle error", "topic", m.Topic, "err", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			c.log.Error("kafkaConsumer.Start: commit error", "topic", m.Topic, "err", err)
			continue
		}

		c.log.Debug("kafkaComsumer.Start: commited")
	}
}

func (c *Consumer) Stop(ctx context.Context) error {
	if err := c.reader.Close(); err != nil {
		c.log.Error("kafkaConsumer.Close: kafka close error", "err", err)
		return err
	}
	return nil
}

func (c *Consumer) process(ctx context.Context, payload []byte) error {
	var msg DTOOrder

	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("Ошибка десериализации json: %w", err)
	}

	order := toDomainOrder(&msg)

	_, err := c.repo.UpsertOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("Ошибка при сохранении заказа в базу данных. %w", err)
	}

	c.cache.Set(order.OrderUID, order)
	c.log.Info("kafkaConsumer.process: processed", "order_uid", order.OrderUID)

	return nil
}
