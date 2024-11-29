.PHONY: test-coverage build-all run-all lint

build-all: ### Полная сборка. Используйте CACHE=1 для кэширования
	@if [ -z "$(CACHE)" ]; then \
		cd cart && make build; \
	else \
		cd cart && make build-cache; \
	fi
	@if [ -z "$(CACHE)" ]; then \
		cd loms && make build; \
	else \
		cd loms && make build-cache; \
	fi
	cd notifier && make build



run-all: build-all ### Поднять и запустить сервисы в docker
	docker-compose up --force-recreate --build -d --scale notifier=3 notifier loms cart
	sleep 3;
	cd loms && make apply-migrations

test-all: ### Запустить тесты на всех сервисах
	cd cart && make test
	cd loms && make test

run-postgres: ### Поднять postgres инфраструктуру из шардов
	docker-compose up -d postgres_slave_0 postgres_master_1

run-observ-infra: ### Поднять инфраструктуру для observability
	docker-compose up -d prometheus grafana jaeger

make clean-postgres: ### Остановить и удалить контейнеры postgres
	docker-compose down postgres_master_0 postgres_master_1 postgres_slave -v
	docker-compose up -d postgres_slave
	sleep 2;
	cd loms && make apply-migrations
	docker-compose down postgres_master postgres_slave
run-kafka: ### Поднять postgres_slave
	docker-compose up -d kafka kafka-ui kafka-init


integration-test-all: ### Запустить интеграционные тесты
	cd loms && go test -v -tags=e2e ./internal/tests/e2e/...
	#cd cart && go test -v -tags=e2e ./internal/tests/e2e/...
# Документация
help: ## Показать этот справочник
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

