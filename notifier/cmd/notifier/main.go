package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"route256/notifier/internal/app/initialization"
	"sync"
	"syscall"
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

	// Создаем контекст с отменой по сигналу
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // Освобождаем ресурсы после завершения

	wg := &sync.WaitGroup{}

	// Запускаем consumer group с переданным контекстом и обработкой graceful shutdown
	wg.Add(1)
	go func() {
		defer wg.Done()
		application.ConsumerGroup.Run(ctx, wg)
	}()

	// Ожидаем завершения работы по сигналу
	<-ctx.Done()
	log.Println("[main] Shutting down Consumer...")

	// Закрываем consumer group
	if err := application.ConsumerGroup.Close(); err != nil {
		log.Printf("[main] Error closing consumer group: %v", err)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	log.Println("[main] Consumer gracefully stopped")
}
