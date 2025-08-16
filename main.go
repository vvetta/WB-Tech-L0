package main

import (
	"time"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

	//TODO Инициализация базы данных.
	err = database.InitDB()
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
