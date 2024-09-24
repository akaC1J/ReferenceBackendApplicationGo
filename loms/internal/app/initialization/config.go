package initialization

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	StockFilePath string
	GgrpcHostPort string
	HttpPort      int
	SwagerUrl     string
}

const defaultHostPortGrpc = ":50051"
const defaultHttpPort = 8081

func LoadDefaultConfig() (*Config, error) {
	return LoadConfig("./.env")
}

func LoadConfig(pathToEnv string) (*Config, error) {
	if err := loadEnv(pathToEnv); err != nil {
		return nil, fmt.Errorf("[config] failed to load environment variables: %w", err)
	}

	stockFilePath := os.Getenv("STOCK_FILE_PATH")

	grpcPort := os.Getenv("GRPC_HOST_PORT")
	if grpcPort == "" {
		log.Printf("[config] will be using default host port: %s", defaultHostPortGrpc)
		grpcPort = defaultHostPortGrpc
	}
	httpPortStr := os.Getenv("HTTP_PORT")
	var httpPort = defaultHttpPort

	if httpPortStr != "" {
		var err error
		httpPort, err = strconv.Atoi(httpPortStr)
		if err != nil {
			log.Printf("[config] failed to parse HTTP_PORT: %v", err)
			log.Printf("[config] will be using default port: %d", defaultHttpPort)
		}
	} else {
		log.Printf("[config] will be using default port: %d", defaultHttpPort)
	}

	swaggerUrl := os.Getenv("SWAGGER_FOR_CORS_ALLOWED_URL")

	return &Config{
		StockFilePath: stockFilePath,
		GgrpcHostPort: grpcPort,
		HttpPort:      httpPort,
		SwagerUrl:     swaggerUrl,
	}, nil
}

func loadEnv(pathToEnv string) error {
	if err := godotenv.Load(pathToEnv); err != nil {
		return fmt.Errorf("failed to load %s: %w", pathToEnv, err)
	}
	return nil
}
