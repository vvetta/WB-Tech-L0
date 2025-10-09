package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type OrderMessage struct {
	OrderUID        string      `json:"order_uid"`
	TrackNumber     string      `json:"track_number"`
	Entry           string      `json:"entry"`
	Delivery        DeliveryMsg `json:"delivery"`
	Payment         PaymentMsg  `json:"payment"`
	Items           []ItemMsg   `json:"items"`
	Locale          string      `json:"locale"`
	InternalSign    string      `json:"internal_signature"`
	CustomerID      string      `json:"customer_id"`
	DeliveryService string      `json:"delivery_service"`
	ShardKey        string      `json:"shardkey"`
	SMID            int         `json:"sm_id"`
	DateCreated     string      `json:"date_created"`
	OOFShard        string      `json:"oof_shard"`
}

type DeliveryMsg struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}
type PaymentMsg struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}
type ItemMsg struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func genOrder(i int) OrderMessage {
	now := time.Now().UTC().Format(time.RFC3339)
	uid := fmt.Sprintf("demo-%d-%d", time.Now().Unix(), i)

	return OrderMessage{
		OrderUID:    uid,
		TrackNumber: "WB123",
		Entry:       "WBIL",
		Delivery: DeliveryMsg{
			Name:    "test name",
			Phone:   "+7 777 77 77",
			Zip:     "test zip",
			City:    "WB City",
			Address: "Pervaya wb ulitsa",
			Region:  "wb region",
			Email:   "nikita@wb.ru",
		},
		Payment: PaymentMsg{
			Transaction:  "hello tr",
			RequestID:    "hello id",
			Currency:     "RU",
			Provider:     "pro test",
			Amount:       1000,
			PaymentDT:    1,
			Bank:         "WB bank",
			DeliveryCost: 20,
			GoodsTotal:   5,
			CustomFee:    8,
		},
		Items: []ItemMsg{
			{ChrtID: 1,
				TrackNumber: "",
				Price:       12,
				RID:         "",
				Name:        "",
				Sale:        12,
				Size:        "Bolshoi",
				TotalPrice:  228,
				NmID:        5,
				Brand:       "wb",
				Status:      1,
			},
		},
		Locale:          "ru",
		InternalSign:    "",
		CustomerID:      "demo",
		DeliveryService: "nikita",
		ShardKey:        "1",
		SMID:            99,
		DateCreated:     now,
		OOFShard:        "1",
	}
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("json marshal: %v", err)
	}
	return b
}

func main() {
	var (
		brokers = flag.String("brokers", "localhost:9092", "comme-separated Kafka brokers")
		topic   = flag.String("topic", "orders", "Kafka topic")
		genN    = flag.Int("gen", 0, "generate N demo orders (if > 0)")
	)

	flag.Parse()

	var payloads [][]byte
	switch {
	case *genN > 0:
		payloads = make([][]byte, 0, *genN)
		for i := 0; i < *genN; i++ {
			payloads = append(payloads, mustJSON(genOrder(i)))
		}
	default:
		log.Fatalf("no input: use --gen")
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(strings.Split(*brokers, ",")...),
		Topic:        *topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		Compression:  kafka.Snappy,
	}
	defer w.Close()

	ctx := context.Background()

	sendBatch := func() error {
		msgs := make([]kafka.Message, 0, len(payloads))
		for _, p := range payloads {
			m := kafka.Message{Value: p}
			msgs = append(msgs, m)
		}
		return w.WriteMessages(ctx, msgs...)
	}

	if err := sendBatch(); err != nil {
		log.Fatalf("kafka write: %v", err)
	}
}
