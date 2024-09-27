package initialization

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	StockFilePath   string
	GgrpcHostPort   string
	HttpPort        int
	SwagerUrl       string
	DBConfigMaster  *DBConfig
	DBConfigReplica *DBConfig
}

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHostPort string
	DBName     string
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
	configMaster, configReplica := MustLoadDBConfig()
	return &Config{
		StockFilePath:   stockFilePath,
		GgrpcHostPort:   grpcPort,
		HttpPort:        httpPort,
		SwagerUrl:       swaggerUrl,
		DBConfigMaster:  configMaster,
		DBConfigReplica: configReplica,
	}, nil
}

func MustLoadDBConfig() (masterConfig, replicaConfig *DBConfig) {
	masterConfig = &DBConfig{
		DBUser:     os.Getenv("DATABASE_MASTER_USER"),
		DBPassword: os.Getenv("DATABASE_MASTER_PASSWORD"),
		DBHostPort: os.Getenv("DATABASE_MASTER_HOST_PORT"),
		DBName:     os.Getenv("DATABASE_MASTER_NAME"),
	}
	switch "" {
	case masterConfig.DBUser:
		log.Fatalf("DATABASE_MASTER_USER is not set")
	case masterConfig.DBPassword:
		log.Fatalf("DATABASE_MASTER_PASSWORD is not set")
	case masterConfig.DBHostPort:
		log.Fatalf("DATABASE_MASTER_HOST_PORT is not set")
	case masterConfig.DBName:
		log.Fatalf("DATABASE_MASTER_NAME is not set")
	}

	replicaConfig = &DBConfig{
		DBUser:     os.Getenv("DATABASE_REPLICA_USER"),
		DBPassword: os.Getenv("DATABASE_REPLICA_PASSWORD"),
		DBHostPort: os.Getenv("DATABASE_REPLICA_HOST_PORT"),
		DBName:     os.Getenv("DATABASE_REPLICA_NAME"),
	}

	//если чего-то не хватает делаем nil, значит реплики у нас не будет
	switch "" {
	case masterConfig.DBUser:
		fallthrough
	case masterConfig.DBPassword:
		fallthrough
	case masterConfig.DBHostPort:
		fallthrough
	case masterConfig.DBName:
		replicaConfig = nil
	}

	return
}

func loadEnv(pathToEnv string) error {
	if err := godotenv.Load(pathToEnv); err != nil {
		return fmt.Errorf("failed to load %s: %w", pathToEnv, err)
	}
	return nil
}
