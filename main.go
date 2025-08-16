package main

import (
	"fmt"
	"log"
	"time"
	"strings"
	"net/http"
	"encoding/json"

	"WB-Tech-L0/cache"
	"WB-Tech-L0/config"
	"WB-Tech-L0/models"
	"WB-Tech-L0/producer"
	"WB-Tech-L0/consumer"
	"WB-Tech-L0/database"
)

func main() {
	err := config.Init()
	if err != nil {
		return
	}

	err = producer.AddTestMessageToKafka()
	if err != nil {
		return
	}

	err = database.InitDB()
	if err != nil {
		return
	}

	err = WarmUpCache()
	if err != nil {
		return
	}

	go consumer.ConsumeOrders()

	http.HandleFunc("/order/", ProcessingOrder)

	http.ListenAndServe(":8081", nil)
}

func ProcessingOrder(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(req.URL.Path, "/order/")
	if id == "" {
		http.NotFound(w, req)
		return
	}

	cachedOrder := cache.Get(id)
	if cachedOrder != nil {
		log.Printf("Заказ был получен из кэша: %s", id)
		log.Printf("Время выполнения запроса: %s", time.Since(startTime))	
		json.NewEncoder(w).Encode(cachedOrder)
	}

	var orders []models.Order
	orders, err := database.GetOrdersFromDB([]string{id})
	if err != nil {
		return
	}

	cache.Set(&orders[0])

	log.Printf("Заказ был получен из базы данных: %s", id)
	log.Printf("Время выполнения запроса: %s", time.Since(startTime))	

	json.NewEncoder(w).Encode(orders[0])
}

func WarmUpCache() error {
// Заполняет кеш данными при запуске сервиса.
	
	var err error
	var orders []models.Order

	orders, err = database.GetOrdersFromDB([]string{}, "all")
	if err != nil {
		return fmt.Errorf("Не получилось загрузить кеш! %w", err)
	}

	for i := 0; i < len(orders) ; i++ {
		cache.Set(&orders[i])
	}
	
	return nil
}

