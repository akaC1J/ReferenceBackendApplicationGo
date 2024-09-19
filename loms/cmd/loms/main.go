package main

import (
	"fmt"
	"log"
	"net"
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

	application, err := app.New(config)
	if err != nil {
		log.Fatalf("[main] Failed to initialize application: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		panic(err)
	}

	log.Println("[main] Application initialization successful")
	log.Printf("[main] server listening at %v", lis.Addr())

	go func() {
		if err = application.GrpcServer.Serve(lis); err != nil {
			log.Fatalf("[main] failed to serve: %v", err)
		}
	}()
	log.Printf("[main] server listening at %v", application.GwServer.Addr)
	if err = application.GwServer.ListenAndServe(); err != nil {
		log.Fatalf("[main] failed to serve: %v", err)
	}

}
