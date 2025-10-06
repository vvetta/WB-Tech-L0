# ======== MIGRATIONS via dockerized golang-migrate ========

# Подтянуть переменные из .env, если файл есть (не обязательно)
-include .env
export

# Где лежат миграции
MIGRATIONS_DIR ?= ./migrations

# DSN БД: можно задать PG_DSN целиком или по частям (PG_USER/PG_PASSWORD/PG_HOST/PG_PORT/PG_DB)
# Пример PG_DSN:
#   postgres://postgres:postgres@localhost:5432/wb_l0?sslmode=disable
DB_DSN ?= $(if $(PG_DSN_URL),$(PG_DSN_URL),postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable)

# Образ migrate
MIGRATE_IMAGE ?= migrate/migrate:latest
MIGRATE_RUN = docker run --rm -v $(PWD)/migrations:/migrations --network host $(MIGRATE_IMAGE)

.PHONY: test cover coverhtml help migrate-up migrate-down migrate-down-all migrate-version migrate-force migrate-new

help:
	@echo "Targets:"
	@echo "  make migrate-up             # применить все неприменённые миграции"
	@echo "  make migrate-down           # откатить одну миграцию"
	@echo "  make migrate-down-all       # откатить все миграции до нуля"
	@echo "  make migrate-version        # показать текущую версию миграций"
	@echo "  make migrate-force v=NNNN   # принудительно выставить версию (аккуратно!)"
	@echo "  make migrate-new name=xxxx  # создать пустые up/down файлы миграции"
	@echo "  make test  # запускает тесты"
	@echo "  make cover  # показывает покрытие кода тестами"
	@echo "  make coverhtml  # выводит html файл с покрытием"

# Применить все недостающие миграции
migrate-up:
	$(MIGRATE_RUN) -path=/migrations -database "$(DB_DSN)" up

# Откатить одну миграцию
migrate-down:
	$(MIGRATE_RUN) -path=/migrations -database "$(DB_DSN)" down 1

# Откатить все миграции
migrate-down-all:
	$(MIGRATE_RUN) -path=/migrations -database "$(DB_DSN)" down

# Показать текущую версию
migrate-version:
	$(MIGRATE_RUN) -path=/migrations -database "$(DB_DSN)" version

# Принудительно установить версию (на случай битого состояния)
migrate-force:
	@if [ -z "$(v)" ]; then echo "Usage: make migrate-force v=<version>"; exit 1; fi
	$(MIGRATE_RUN) -path=/migrations -database "$(DB_DSN)" force $(v)

# Создать заготовки up/down с таймстемпом и именем
migrate-new:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-new name=<short_name>"; exit 1; fi ;\
	ts=$$(date +%Y%m%d%H%M%S); \
	touch $(MIGRATIONS_DIR)/$${ts}_$(name).up.sql; \
	touch $(MIGRATIONS_DIR)/$${ts}_$(name).down.sql; \
	echo "Created: $(MIGRATIONS_DIR)/$${ts}_$(name).up.sql / .down.sql"

test:
	go test ./... -cover

cover:
	go test ./... -covermode=atomic -coverprofile=coverage.out && \
	go tool cover -func=coverage.out

coverhtml:
	go test ./... -covermode=atomic -coverprofile=coverage.out && \
	go tool cover -html=coverage.out
