# Переменные
BIN_DIR := $(PWD)/bin
REPORTS_DIR := $(PWD)/reports
MINIMOCK_URL := github.com/gojuno/minimock/v3/cmd/minimock@latest
BUF_TOOL_URL := github.com/bufbuild/buf/cmd/buf@v1.41.0
.PHONY: build test install-tools generate_mocks help

export GOBIN := $(BIN_DIR)

# Цель для сборки
build: ## Собрать проект
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o  $(BIN_DIR)/loms cmd/loms/main.go

#Цель для локального запуска
fast-run: ## Запустить проект
	go run cmd/loms/main.go

#Цель для протокомпиляции
buf: install-tools ## кодогенерация протофайлов
	$(BIN_DIR)/buf generate
	go mod tidy

# Цель для установки инструментов
install-tools: ## Установить необходимые инструменты (minimock)
	@mkdir -p $(BIN_DIR)
	@[ -f $(BIN_DIR)/minimock ] || { \
		echo >&2 "Installing minimock..."; \
		go install $(MINIMOCK_URL); \
	}
	@[ -f $(BIN_DIR)/buf ] || { \
    		echo >&2 "Installing buf..."; \
    		go install $(BUF_TOOL_URL); \
    	}

# Команды для тестирования
test: generate_mocks ## Запустить тесты с покрытием
	@mkdir -p $(REPORTS_DIR)
	go test -coverprofile=$(REPORTS_DIR)/coverage.out ./... && \
	go tool cover -html=$(REPORTS_DIR)/coverage.out -o $(REPORTS_DIR)/coverage.html

generate_mocks: install-tools ## Генерация моков
	@echo "Generating mocks for all packages"
	$(BIN_DIR)/minimock -i route256/loms/internal/service/orderservice.StockService -o internal/service/orderservice/test
	$(BIN_DIR)/minimock -i route256/loms/internal/service/orderservice.Repository -o internal/service/orderservice/test
	$(BIN_DIR)/minimock -i route256/loms/internal/service/stockService.Repository -o internal/service/stockservice/test

# Документация
help: ## Показать этот справочник
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
