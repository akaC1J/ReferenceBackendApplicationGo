BIN_DIR := $(PWD)/bin



# Цель для сборки
.PHONY: build
build:  ## Собрать проект
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/notifier ./cmd/notifier

#Цель для быстрого запуска
.PHONY: fast-run
fast-run:
	go run cmd/notifier/main.go

# Цель для запуска тестов
.PHONY: test
test: ## Запустить тесты с покрытием
	go test ./...


# Цель для отображения справки
.PHONY: help
help: ## Показать этот справочник
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
