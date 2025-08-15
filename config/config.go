/*
Package config содержит в себе всю конфигурацию приложения.
*/
package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Init() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Ошибка загрузки .env файла: %v", err)
		return err
	}

	log.Print("Конфигурация приложения был успешно загружена!")
	return nil
}

func GetKafkaBrokers() ([]string, error) {
// Получает адреса клиентов Kafka из .env файла.
	hosts := strings.Split(os.Getenv("KAFKAHOSTS"), ",")
	ports := strings.Split(os.Getenv("KAFKAPORTS"), ",")

	if (len(hosts) != len(ports)) || len(hosts) == 0 || len(ports) == 0 {
		log.Print("Не получилось загрузить Kafka адреса!")
		return []string{}, fmt.Errorf("Не получилось загрузить Kafka адреса!")
	}

	var kafkaBrokers []string

	for i := 0; i < len(hosts); i++ {
		broker := string(hosts[i]) + ":" + string(ports[i])
		kafkaBrokers = append(kafkaBrokers, broker)
	}

	return kafkaBrokers, nil
}

func GetKafkaTopicName() string {
// Возвращает имя топика из .env файла, если его нет, то возвращает дефолтное значение.
	topicName := os.Getenv("KAFKATOPICNAME")

	if len(topicName) == 0 {
		topicName = "ordersTopic"
	}

	return topicName
}

func GetDSN() string {
// Возвращает данные для подключения к базе данных. DSN (Data Source Name)
	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
	dbPort := os.Getenv("DBPORT")

	if dbHost == "" {
		log.Fatal("В .env файле отсутствует 'DBHOST'.")
	} else if dbUser == "" {
		log.Fatal("В .env файле отсутствует 'DBUSER'.")
	} else if dbPass == "" {
		log.Fatal("В .env файле отсутствует 'DBPASS'.")
	} else if dbName == "" {
		log.Fatal("В .env файле отсутствует 'DBNAME'.")
	} else if dbPort == "" {
		log.Fatal("В .env файле отсутствует 'DBPORT'.") 
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPass, dbName, dbPort)

	return dsn
}
