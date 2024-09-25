package main

import (
	"log"
	"os"
	app "route256/loms/internal/app/initialization"
)

func main() {
	envPath := os.Getenv("ENV_PATH")
	var config *app.Config
	var err error
	if envPath != "" {
		config, err = app.LoadConfig(envPath)
	} else {
		config, err = app.LoadDefaultConfig()
	}
	if err != nil {
		log.Fatalf("[main] Failed to load configuration: %v", err)
	}

	application, err := app.MustNew(config)
	if err != nil {
		log.Fatalf("[main] Failed to initialize application: %v", err)
	}

	application.Run(config)
}
