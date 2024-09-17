.PHONY: test-coverage build-all run-all lint

build-all:
	cd cart && make build
	cd loms && make build

run-all: build-all
	docker-compose up --force-recreate --build -d

test:
	cd cart && make test
	cd loms && make test
