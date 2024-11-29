# Переменные
BIN_DIR := $(PWD)/bin
REPORTS_DIR := $(PWD)/reports
MINIMOCK_URL := github.com/gojuno/minimock/v3/cmd/minimock@latest
PROTOC_GEN_VALIDATOR_URL := github.com/envoyproxy/protoc-gen-validate@v1.0.4
PROTOC_GEN_GATEWAY_URL := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
PROTOC_GEN_OPENAPI_URL := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
PROTOC_GEN_GO_URL := google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
PROTOC_GEN_GO_GRPC_URL := google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
GOOSE_URL := github.com/pressly/goose/v3/cmd/goose@v3.22.1
SQLC_URL := github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

export GOBIN := $(BIN_DIR)

.PHONY: build
# Цель для сборки

build: .install-tools .protoc-generate .sqlc-generate ### Собрать проект
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o  $(BIN_DIR)/loms cmd/loms/main.go

build-cache:  ## Собрать проект не устанавливая зависимости
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/loms ./cmd/loms

.PHONY: fast-run
#Цель для локального запуска
fast-run: ## Запустить проект
	go run cmd/loms/main.go

.PHONY: test
# Команды для тестирования
test: .generate_mocks ## Запустить тесты с покрытием
	@mkdir -p $(REPORTS_DIR)
	go test -coverprofile=$(REPORTS_DIR)/coverage.out ./... && \
	go tool cover -html=$(REPORTS_DIR)/coverage.out -o $(REPORTS_DIR)/coverage.html

.PHONY: .generate_mocks
.generate_mocks: .install-tools ## Генерация моков
	@echo "Generating mocks for all packages"
	$(BIN_DIR)/minimock -i route256/loms/internal/service/orderservice.StockService -o internal/service/orderservice/test
	$(BIN_DIR)/minimock -i route256/loms/internal/service/orderservice.Repository -o internal/service/orderservice/test
	$(BIN_DIR)/minimock -i route256/loms/internal/service/stockService.Repository -o internal/service/stockservice/test

.PHONY: help
# Документация
help: ## Показать этот справочник
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
LOMS_PROTO_PATH:="api/loms/v1"
.PHONY: .protoc-generate
.protoc-generate: ## Генерация кода из proto-файлов
	@mkdir -p ./api/openapiv2
	protoc \
	-I ${LOMS_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go \
	--go_out internal/generated/${LOMS_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc \
	--go-grpc_out internal/generated/${LOMS_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(BIN_DIR)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:internal/generated/api/loms/v1" \
	--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway \
	--grpc-gateway_out internal/generated/${LOMS_PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2 \
	--openapiv2_out api/openapiv2 \
	--openapiv2_opt logtostderr=true \
	loms.proto
	go mod tidy


.PHONY: .sqlc-generate
sqlc-generate: ## Генерация кода из sql-файлов
	@echo "Generating sqlc code"
	$(BIN_DIR)/sqlc generate

ifneq (,$(wildcard ../.env))
    include ../.env
    export
endif
.PHONY: apply-migrations
apply-migrations: ## Применить миграции goose
	$(foreach idx,0 1,\
		bin/goose -dir migrations/common postgres \
		"postgresql://$(POSTGRES_MASTER_USER_$(idx)):$(POSTGRES_MASTER_PASSWORD_$(idx))@$(POSTGRES_MASTER_HOST_PORT_$(idx))/loms?sslmode=disable" up; \
		bin/goose -dir migrations/shard_$(idx) postgres \
		"postgresql://$(POSTGRES_MASTER_USER_$(idx)):$(POSTGRES_MASTER_PASSWORD_$(idx))@$(POSTGRES_MASTER_HOST_PORT_$(idx))/loms?sslmode=disable" up; \
	)

.PHONY: add_migrations
add_migrations: .install-tools ### Создать новый файл миграции, обязателен параметр MIGR_NAME
ifeq ($(MIGR_NAME),)
	$(error Необходимо указать имя файла через параметр MIGR_NAME)
endif
	@echo "Создание файла: $(MIGR_NAME)"
	@mkdir -p migrations
	@cd migrations && $(BIN_DIR)/goose create $(MIGR_NAME) sql

.PHONY: .install-tools
# Цель для установки инструментов
.install-tools:
	@mkdir -p $(BIN_DIR)
	@[ -f $(BIN_DIR)/minimock ] || { \
		echo >&2 "Installing minimock..."; \
		go install $(MINIMOCK_URL); \
	};
	@[ -f $(BIN_DIR)/protoc-gen-validate ] || { \
    		echo >&2 "Installing protoc-gen-validator..."; \
    		go install $(PROTOC_GEN_VALIDATOR_URL); \
    }
	@[ -f $(BIN_DIR)/protoc-gen-go ] || { \
    		echo >&2 "Installing protoc-gen-go..."; \
    		go install $(PROTOC_GEN_GO_URL); \
    }
	@[ -f $(BIN_DIR)/protoc-gen-go-grpc ] || { \
    		echo >&2 "Installing protoc-gen-go_grpc..."; \
    		go install $(PROTOC_GEN_GO_GRPC_URL); \
    }
	@[ -f $(BIN_DIR)/protoc-gen-grpc-gateway ] || { \
    		echo >&2 "Installing protoc-gen-grpc-gateway..."; \
    		go install $(PROTOC_GEN_GATEWAY_URL); \
    }
	@[ -f $(BIN_DIR)/protoc-gen-openapiv2 ] || { \
    		echo >&2 "Installing protoc-gen-openapi..."; \
    		go install $(PROTOC_GEN_OPENAPI_URL); \
    }
	@[ -f $(BIN_DIR)/goose ] || { \
    		echo >&2 "Installing goose..."; \
    		go install $(GOOSE_URL); \
    }
	@[ -f $(BIN_DIR)/sqlc ] || { \
    		echo >&2 "Installing sqlc..."; \
    		go install $(SQLC_URL); \
    }

.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp


# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf
