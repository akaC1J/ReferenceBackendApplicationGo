package main

import (
	"log"
	"net/http"
	"os"
	"route256/cart/internal/app/initialization"
)

func main() {
	envPath := os.Getenv("ENV_PATH")
	var config *initialization.Config
	var err error
	if envPath != "" {
		config, err = initialization.LoadConfig(envPath)
	} else {
		config, err = initialization.LoadDefaultConfig()
	}

	if err != nil {
		log.Fatalf("[main] Failed to load configuration: %v", err)
	}

	application, err := initialization.New(config)
	if err != nil {
		log.Fatalf("[main] Failed to initialize application: %v", err)
	}

	log.Println("[main] Application initialization successful")

	log.Printf("[main] Starting server on %s\n", config.HostPort)
	if err := http.ListenAndServe(config.HostPort, application); err != nil {
		log.Fatalf("[main] Failed to start server: %v", err)
	}
}
