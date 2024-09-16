.PHONY: test-coverage build-all run-all lint

build-all:
	cd cart && make build

run-all: build-all
	docker-compose up --force-recreate --build -d

test-coverage:
	cd cart && make test-coverage

lint:
	cd cart && make lint
