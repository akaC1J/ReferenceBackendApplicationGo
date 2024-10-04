package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"route256/cart/internal/app/initialization"
	"syscall"
	"time"
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

	srv := &http.Server{
		Addr:    config.HostPort,
		Handler: application,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[main] Starting server on %s\n", config.HostPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("[main] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[main] Server forced to shutdown: %v", err)
	}

	log.Println("[main] Server gracefully stopped")
}
