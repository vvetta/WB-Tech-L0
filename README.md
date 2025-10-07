# WB-Tech-L0 - Order Service

Микросервис для приема, хранения и отображения информации о заказах.

Сервис получает JSON-сообщения из **Kafka**, сохраняет данные о заказах в **PostgreSQL**, кэширует их в памяти, также предоставляет HTTP-API и простой web-интерфейс для просмотра заказа по `order_uid`.

## Архитектура проекта

Проект реализован в соответствии с принципами **Чистой архитектуры**
(разделение по слоям: `domain` -> `usecase` -> `adapters` -> `cmd`).

## Используемы технологии

|Компонент|Технология / библиотека|
|---------|-----------------------|
|Язык     |Go 1.22+               |
|Бд       |PostgreSQL + GORM      |
|Очередь  |Kafka (segmentio/kafka-go)|
|HTTP-сервер|net/http             |
|Кеш      |map + sync.Mutex       |
|Миграции |golang-migrate + Makefile|
|Логирование|slog|
|Генерация моков|golang/mock      |
|Тесты    |testify + gomock       |
|Фейковые данные|gofakeit         |
|Frontend |HTML + JS + live-server (for dev)|

## Переменные окружения

Перед запуском создайте файл `.env` в корне проекта:
Пример `.env` файла описан в файле `.env-example`

## Локальный запуск

### 1. Установить зависимости

`go mod download`

### 2. Поднять инфраструктуру

Необходимо запустить PostgreSQL и Kafka.

- [Инструкция по запуску Kafka в Docker](https://purpleschool.ru/knowledge-base/article/kafka)
- [Инструкция по запуску PostgreSQL в Docker](https://habr.com/ru/articles/578744/)

### 3. Применить миграции

`make migrate-up`

### 4. Запустить backend

`go run ./cmd/main.go`

### 5. Запустить frontend

`npm install -g live-server`
`cd frontend`
`live-server --port=3000`

## Тестирование

Для запуска тестов используйте команду: `make test`
Для генерации HTML-отчета о покрытии: `make coverhtml`

## Пример JSON-сообщения заказа

```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}

```
```

```
