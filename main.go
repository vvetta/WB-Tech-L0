package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(req.URL.Path, "/order/")
	if id == "" {
		http.NotFound(w, req)
		return
	}

	//TODO Получение заказа из базы данных.

	//TODO Кэширование заказа.

	log.Printf("Начинаю поиск заказа: %s", id)

	resp := models.Order{
		OrderUID: id,
	}

	json.NewEncoder(w).Encode(resp)
}
