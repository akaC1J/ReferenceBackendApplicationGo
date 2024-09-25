.PHONY: test-coverage build-all run-all lint

build-all: ### Полная сборка
	cd cart && make build
	cd loms && make build

run-all: build-all run-postgres ### Поднять и запустить сервисы в docker
	docker-compose up --force-recreate --build -d

test-all: ### Запустить тесты на всех сервисах
	cd cart && make test
	cd loms && make test

run-postgres: ### Поднять postgres
	docker-compose up -d postgres

integration-test-all: ### Запустить интеграционные тесты
	cd loms && go test -v -tags=e2e ./internal/tests/e2e/...
	cd cart && go test -v -tags=e2e ./internal/tests/e2e/...
# Документация
help: ## Показать этот справочник
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

