build-all:
	cd cart && make build


run-all: build-all
	docker-compose up --force-recreate --build -d
