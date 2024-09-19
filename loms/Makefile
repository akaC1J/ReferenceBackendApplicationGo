# Переменные
BIN_DIR := $(PWD)/bin
REPORTS_DIR := $(PWD)/reports
MINIMOCK_URL := github.com/gojuno/minimock/v3/cmd/minimock@latest
PROTOC_GEN_VALIDATOR_URL := github.com/envoyproxy/protoc-gen-validate@v1.0.4
PROTOC_GEN_GATEWAY_URL := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
PROTOC_GEN_OPENAPI_URL := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
SWAGGER_URL := github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5
PROTOC_GEN_GO_URL := google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
PROTOC_GEN_GO_GRPC_URL := google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

.PHONY: build test install-tools generate_mocks help .vendor-proto protoc-generate

export GOBIN := $(BIN_DIR)

# Цель для сборки
build: install-tools .protoc-generate ## Собрать проект
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o  $(BIN_DIR)/loms cmd/loms/main.go

#Цель для локального запуска
fast-run: ## Запустить проект
	go run cmd/loms/main.go

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

LOMS_PROTO_PATH:="api/loms/v1"
.protoc-generate: ## Генерация кода из proto-файлов
	@mkdir -p ./api/openapiv2
	protoc \
	-I ${LOMS_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go \
	--go_out pkg/${LOMS_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc \
	--go-grpc_out pkg/${LOMS_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(BIN_DIR)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:pkg/api/loms/v1" \
	--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway \
	--grpc-gateway_out pkg/${LOMS_PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2 \
	--openapiv2_out api/openapiv2 \
	--openapiv2_opt logtostderr=true \
	api/loms/v1/loms.proto
	go mod tidy

.vendor-proto: .vendor-rm install-tools  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

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


# Цель для установки инструментов
install-tools:
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
