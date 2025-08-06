/*
Package config содержит в себе всю конфигурацию приложения.
*/
package config

import (
	"log"

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
