package main

import (
	"log"
	"net/http"

	"route256/cart/internal/app"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("[main] Failed to load configuration: %v", err)
	}

	application, err := app.New(config)
	if err != nil {
		log.Fatalf("[main] Failed to initialize application: %v", err)
	}

	log.Println("[main] Application initialization successful")

	log.Printf("[main] Starting server on %s\n", config.Host_Port)
	if err := http.ListenAndServe(config.Host_Port, application); err != nil {
		log.Fatalf("[main] Failed to start server: %v", err)
	}
}
