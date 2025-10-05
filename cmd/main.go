package main

import (
	"context"
	"strings"
	"fmt"
	"time"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"

	"WB-Tech-L0/internal/adapters/http"
	"WB-Tech-L0/internal/adapters/cache"
	"WB-Tech-L0/internal/adapters/kafka"
	"WB-Tech-L0/internal/adapters/log"
	"WB-Tech-L0/internal/adapters/repository"
	"WB-Tech-L0/internal/usecase"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustAtoi(s string, def int) int {
	if s == "" {
		return def
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return def
	}
	return n
}

func main() {
	_ = godotenv.Load(".env")

	pgDSN := getenv("PG_DSN", "host=localhost port=5432 user=postgres password=postgres dbname=wb_l0")

	cacheLimit := mustAtoi(os.Getenv("CACHE_LIMIT"), 10)
	cacheDeleteCount := mustAtoi(os.Getenv("CACHE_DELETE_COUNT"), 2)
	warmUpLimit := cacheLimit - 1

	kafkaBrokers := strings.Split(getenv("KAFKA_BROKERS", "localhost:9092"), ",")
	kafkaTopic := getenv("KAFKA_TOPIC", "order")
	kafkaGroup := getenv("KAFKA_GROUP_ID", "order-service")

	httpAddr := getenv("HTTP_ADDR", ":8081")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lg := logger.NewLogger()
	lg.Info("boot: starting service")

	db, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
	if err != nil {
		lg.Error("boot: gorm open failed", "err", err)
		log.Fatalf("gorm open failed %v", err)			
	}

	repo := repository.NewPostgresRepo(db, lg)
	memcache := cache.NewMemoryCache(cacheLimit, cacheDeleteCount, lg)
	svc :=	usecase.NewOrderService(repo, memcache, lg) 

	if err := svc.WarmUpCache(ctx, warmUpLimit); err != nil {
		lg.Error("boot: warmUp failed", "err", err)
	}

	consCfg := kafka.Config{
		Brokers: kafkaBrokers,
		Topic: kafkaTopic,
		GroupID: kafkaGroup,
		MinBytes: 1 << 10,
		MaxBytes: 10 << 20,
		StartOffset: -1,
		CommitInterval: 0,
	}
	consumer, err := kafka.New(consCfg, repo, memcache, lg)
	if err != nil {
		lg.Error("boot: kafka init failed", "err", err)
		log.Fatalf("kafka init failed %v", err)	
	}
	go func() {
		if err := consumer.Start(ctx); err != nil {
			lg.Error("kafka: stopped with error", "err", err)
		}
	}()

	srv := http.NewServer(httpAddr, http.Deps{
		OrderSvc: svc,
		Logger: lg,
	})

	go func() {
		lg.Info("http: server starting", "addr", httpAddr)
		if err := srv.Start(); err != nil {
			lg.Error("http: server error", "err", err)
			cancel()	
		}
	}()

	<-ctx.Done()
	lg.Info("boot: shutting down...")

	shCtx, shCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shCancel()

	_ = consumer.Stop(ctx)
	_ = srv.Shutdown(shCtx)
	
	lg.Info("bye!!!")
}
