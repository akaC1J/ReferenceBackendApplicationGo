package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	MaxRetries         int
	RetryDelayMs       int
	Token              string
	ProductServiceURL  string
	ProductServicePath string
	Host_Port          string
	Port               string
}

func LoadDefaultConfig() (*Config, error) {
	return LoadConfig("./.env")
}
func LoadConfig(pathToEnv string) (*Config, error) {
	if err := loadEnv(pathToEnv); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	maxRetries, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse MAX_RETRIES: %w", err)
	}

	retryMs, err := strconv.Atoi(os.Getenv("RETRY_DELAY_MS"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RETRY_DELAY_MS: %w", err)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TOKEN environment variable is required")
	}

	baseUrl := os.Getenv("PRODUCT_SERVICE_BASE_URL")
	if baseUrl == "" {
		return nil, fmt.Errorf("PRODUCT_SERVICE_BASE_URL environment variable is required")
	}

	path := os.Getenv("PRODUCT_SERVICE_PATH")
	if path == "" {
		return nil, fmt.Errorf("PRODUCT_SERVICE_PATH environment variable is required")
	}

	host_port := os.Getenv("HOST_PORT")
	if host_port == "" {
		return nil, fmt.Errorf("HOST environment variable is required")
	}

	return &Config{
		MaxRetries:         maxRetries,
		RetryDelayMs:       retryMs,
		Token:              token,
		ProductServiceURL:  baseUrl,
		ProductServicePath: path,
		Host_Port:          host_port,
	}, nil
}

func loadEnv(pathToEnv string) error {
	if err := godotenv.Load(pathToEnv); err != nil {
		return fmt.Errorf("failed to load %s: %w", pathToEnv, err)
	}
	return nil
}
