# Переменные
BIN_DIR := $(PWD)/bin
REPORTS_DIR := $(PWD)/reports
GOCYCLO_URL := github.com/fzipp/gocyclo/cmd/gocyclo@latest
GOCOGNIT_URL := github.com/uudashr/gocognit/cmd/gocognit@latest
MINIMOCK_URL := github.com/gojuno/minimock/v3/cmd/minimock@latest

# Экспорт переменных окружения
export GOBIN := $(BIN_DIR)

# Определение phony целей
.PHONY: build test test-coverage install-tools cyclo cognit lint generate_mocks help

# Цель для сборки
build: ## Собрать проект
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/cart ./cmd/cart

#Цель для быстрого запуска
fast-run:
	go run cmd/cart/main.go

# Цель для запуска тестов
test: generate_mocks ## Запустить тесты с покрытием
	@mkdir -p $(REPORTS_DIR)
	go test -coverprofile=$(REPORTS_DIR)/coverage.out ./... && \
	go tool cover -html=$(REPORTS_DIR)/coverage.out -o $(REPORTS_DIR)/coverage.html

# Цель для установки инструментов
install-tools: ## Установить необходимые инструменты (gocyclo, gocognit, minimock)
	@mkdir -p $(BIN_DIR)
	@[ -f $(BIN_DIR)/gocyclo ] || { \
		echo >&2 "Installing gocyclo..."; \
		go install $(GOCYCLO_URL); \
	}
	@[ -f $(BIN_DIR)/gocognit ] || { \
		echo >&2 "Installing gocognit..."; \
		go install $(GOCOGNIT_URL); \
	}
	@[ -f $(BIN_DIR)/minimock ] || { \
		echo >&2 "Installing minimock..."; \
		go install $(MINIMOCK_URL); \
	}

# Цель для анализа цикломатической сложности
cyclo: install-tools ## Запустить gocyclo для анализа цикломатической сложности
	@mkdir -p $(REPORTS_DIR)
	$(BIN_DIR)/gocyclo -ignore ".*mock_test.go" . > $(REPORTS_DIR)/cyclo-report.txt

# Цель для анализа когнитивной сложности
cognit: install-tools ## Запустить gocognit для анализа когнитивной сложности
	@mkdir -p $(REPORTS_DIR)
	$(BIN_DIR)/gocognit -ignore ".*mock_test.go" . > $(REPORTS_DIR)/cognit-report.txt

# Цель для линтинга кода
lint: cyclo cognit ## Провести линтинг кода с использованием gocyclo и gocognit

# Цель для генерации моков
.PHONY: generate_mocks
generate_mocks: install-tools ## Генерация моков
	@echo "Generating mocks for all packages"
	$(BIN_DIR)/minimock -i  route256/cart/internal/pkg/repository.* -o internal/pkg/repository
	$(BIN_DIR)/minimock -i route256/cart/internal/pkg/service/cartservice.* -o internal/pkg/service/cartservice
	$(BIN_DIR)/minimock -i route256/cart/internal/pkg/service/productservice.* -o internal/pkg/service/productservice
	$(BIN_DIR)/minimock -i route256/cart/internal/pkg/service/lomsservice.* -o internal/pkg/service/lomsservice
	$(BIN_DIR)/minimock -i route256/loms/pkg/api/loms/v1.LomsClient -o internal/pkg/service/lomsservice

# Цель для отображения справки
help: ## Показать этот справочник
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
