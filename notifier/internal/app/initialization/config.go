package initialization

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupId string
}

type Config struct {
	KafkaConfig KafkaConfig
}

func LoadDefaultConfig() (*Config, error) {
	return LoadConfig("./.env")
}
func LoadConfig(pathToEnv string) (*Config, error) {
	if err := loadEnv(pathToEnv); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	return &Config{
		KafkaConfig: KafkaConfig{
			Brokers: []string{os.Getenv("KAFKA_BROKER")},
			Topic:   os.Getenv("KAFKA_TOPIC"),
			GroupId: os.Getenv("KAFKA_GROUP_ID"),
		},
	}, nil
}

func loadEnv(pathToEnv string) error {
	if err := godotenv.Load(pathToEnv); err != nil {
		return fmt.Errorf("failed to load %s: %w", pathToEnv, err)
	}
	return nil
}
