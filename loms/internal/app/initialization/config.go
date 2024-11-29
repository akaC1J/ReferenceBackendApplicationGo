package initialization

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type Config struct {
	StockFilePath  string
	GgrpcHostPort  string
	HttpPort       int
	SwagerUrl      string
	DBConfigs      []*DBConfigs
	KafkaConfig    *KafkaConfig
	IntervalOutbox time.Duration
}

type DBConfigs struct {
	Master          *DBConfig
	ReplicaOptional *DBConfig
}

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBHostPort string
	DBName     string
}

type KafkaConfig struct {
	Brokers  []string
	RetryMax uint
	Topic    string
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
	dbConfigs, err := MustLoadDBConfig()
	if err != nil {
		log.Printf("[config] replica configuration is incomplete, using only master")
	}

	intervalOutboxStr := os.Getenv("INTERVAL_OUTBOX")
	intervalOutbox, err := strconv.ParseUint(intervalOutboxStr, 10, 64)

	if err != nil || intervalOutbox > math.MaxInt64 {
		log.Printf("[config] failed to parse INTERVAL_OUTBOX: %v", err)
		intervalOutbox = 1_000_000_000
	}
	return &Config{
		StockFilePath:  stockFilePath,
		GgrpcHostPort:  grpcPort,
		HttpPort:       httpPort,
		SwagerUrl:      swaggerUrl,
		DBConfigs:      dbConfigs,
		KafkaConfig:    MustLoadKafkaConfig(),
		IntervalOutbox: time.Duration(intervalOutbox),
	}, nil
}

func MustLoadDBConfig() ([]*DBConfigs, error) {
	var dbConfigs []*DBConfigs
	var err error
	generateEnvNameWithIndex := func(envName string, index int) string {
		return fmt.Sprintf("%s_%d", envName, index)
	}
	for i := 0; i < 2; i++ {
		masterConfig := &DBConfig{
			DBUser:     os.Getenv(generateEnvNameWithIndex("DATABASE_MASTER_USER", i)),
			DBPassword: os.Getenv(generateEnvNameWithIndex("DATABASE_MASTER_PASSWORD", i)),
			DBHostPort: os.Getenv(generateEnvNameWithIndex("DATABASE_MASTER_HOST_PORT", i)),
			DBName:     os.Getenv(generateEnvNameWithIndex("DATABASE_MASTER_NAME", i)),
		}
		switch "" {
		case masterConfig.DBUser:
			log.Fatalf(generateEnvNameWithIndex("DATABASE_MASTER_USER", i) + " is not set")
		case masterConfig.DBPassword:
			log.Fatalf(generateEnvNameWithIndex("DATABASE_MASTER_PASSWORD", i) + " is not set")
		case masterConfig.DBHostPort:
			log.Fatalf(generateEnvNameWithIndex("DATABASE_MASTER_HOST_PORT", i) + " is not set")
		case masterConfig.DBName:
			log.Fatalf(generateEnvNameWithIndex("DATABASE_MASTER_NAME", i) + " is not set")
		}

		replicaConfigOptional := &DBConfig{
			DBUser:     os.Getenv(generateEnvNameWithIndex("DATABASE_REPLICA_USER", i)),
			DBPassword: os.Getenv(generateEnvNameWithIndex("DATABASE_REPLICA_PASSWORD", i)),
			DBHostPort: os.Getenv(generateEnvNameWithIndex("DATABASE_REPLICA_HOST_PORT", i)),
			DBName:     os.Getenv(generateEnvNameWithIndex("DATABASE_REPLICA_NAME", i)),
		}

		//если чего-то не хватает делаем nil, значит реплики у нас не будет
		switch "" {
		case replicaConfigOptional.DBUser:
			fallthrough
		case replicaConfigOptional.DBPassword:
			fallthrough
		case replicaConfigOptional.DBHostPort:
			fallthrough
		case replicaConfigOptional.DBName:
			err = fmt.Errorf("replica with index" + fmt.Sprintf("%d ", i) + "configuration is incomplete")
			replicaConfigOptional = nil
		}
		dbConfigs = append(dbConfigs, &DBConfigs{
			Master:          masterConfig,
			ReplicaOptional: replicaConfigOptional,
		})
	}

	return dbConfigs, err
}

func MustLoadKafkaConfig() (kafkaConfig *KafkaConfig) {
	maxRetryStr := os.Getenv("KAFKA_RETRY_MAX")
	maxRetry, err := strconv.ParseUint(maxRetryStr, 10, 0)
	if err != nil {
		log.Fatalf("KAFKA_RETRY_MAX is incorrect: %v", err)
	}

	kafkaConfig = &KafkaConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		RetryMax: uint(maxRetry),
		Topic:    os.Getenv("KAFKA_TOPIC"),
	}
	if kafkaConfig.Brokers[0] == "" {
		log.Fatalf("KAFKA_BROKER is not set")
	}
	if kafkaConfig.Topic == "" {
		log.Fatalf("KAFKA_TOPIC is not set")
	}
	return kafkaConfig
}

func loadEnv(pathToEnv string) error {
	if err := godotenv.Load(pathToEnv); err != nil {
		return fmt.Errorf("failed to load %s: %w", pathToEnv, err)
	}
	return nil
}
