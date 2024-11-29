package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"route256/cart/internal/app/initialization"
	pkgLogger "route256/cart/internal/logger"
	"syscall"
	"time"
)

func main() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.ErrorOutputPaths = []string{"stdout"}
	zapConfig.Level.SetLevel(zap.InfoLevel)
	zapLogger := pkgLogger.NewLogger(zapConfig)

	envPath := os.Getenv("ENV_PATH")
	var config *initialization.Config
	var err error
	if envPath != "" {
		config, err = initialization.LoadConfig(envPath)
	} else {
		config, err = initialization.LoadDefaultConfig()
	}

	if err != nil {
		pkgLogger.PanicF(nil, "Failed to load configuration: %v\n", err)
	}

	application, err := initialization.New(config, zapLogger)
	if err != nil {
		pkgLogger.PanicF(nil, "Failed to initialize application: %v\n", err)
	}

	pkgLogger.Infow(nil, " Application initialization successful\n")

	srv := &http.Server{
		Addr:    config.HostPort,
		Handler: application,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		pkgLogger.Infow(nil, "[main] Starting server on %s\n", config.HostPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] Failed to start server: %v", err)
		}
	}()

	<-quit
	pkgLogger.Infow(nil, "[main] Shutting down server...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		pkgLogger.PanicF(nil, "[main] Server forced to shutdown: %v\n", err)
	}

	pkgLogger.Infow(nil, "[main] Server gracefully stopped\n")
}
