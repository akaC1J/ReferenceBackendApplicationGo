package initialization

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"route256/cart/internal/logger"
	"strconv"
	"time"
)

type Config struct {
	MaxRetries         int
	RetryDelayMs       int
	Token              string
	ProductServiceURL  string
	ProductServicePath string
	HostPort           string
	Port               string
	LomsBaseUrl        string
	RequestsPerSecond  uint
	CacheCapacity      int
	RedisClient        *RedisConf
}

type RedisConf struct {
	Addr     string
	Password string
	DB_Index int
	KeyTtl   time.Duration
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

	requestsPerSecond, err := strconv.ParseUint(os.Getenv("REQUESTS_PER_SECOND"), 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REQUESTS_PER_SECOND: %w", err)
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

	hostPort := os.Getenv("HOST_PORT")
	if hostPort == "" {
		return nil, fmt.Errorf("HOST environment variable is required")
	}

	lomsBaseUrl := os.Getenv("LOMS_BASE_URL")
	if lomsBaseUrl == "" {
		return nil, fmt.Errorf("LOMS_BASE_URL environment variable is required")
	}

	cacheCapacity := os.Getenv("CACHE_CAPACITY")
	cacheCapacityInt := 10
	if cacheCapacity == "" {
		logger.Warnw(nil, fmt.Sprintf("CACHE_CAPACITY environment variable is not exist. Default value 10 will be used"))
	} else {
		cacheCapacityInt, err = strconv.Atoi(cacheCapacity)
		if err != nil {
			return nil, err
		}
	}

	redisConf, err := initRedisClient()

	if err != nil {
		return nil, err
	}
	return &Config{
		MaxRetries:         maxRetries,
		RetryDelayMs:       retryMs,
		Token:              token,
		ProductServiceURL:  baseUrl,
		ProductServicePath: path,
		HostPort:           hostPort,
		LomsBaseUrl:        lomsBaseUrl,
		RequestsPerSecond:  uint(requestsPerSecond),
		CacheCapacity:      cacheCapacityInt,
		RedisClient:        redisConf,
	}, nil
}

func initRedisClient() (*RedisConf, error) {
	dbIndex := 0
	err := error(nil)
	if os.Getenv("REDIS_DB_INDEX") != "" {
		dbIndex, err = strconv.Atoi(os.Getenv("REDIS_DB_INDEX"))
		if err != nil {
			return nil, fmt.Errorf("enviroment REDIS_DB_INDEX not a number")
		}
	}
	keyTtl := time.Hour
	if os.Getenv("REDIS_TTL_KEY") == "" {
		keyTtlInt, err := strconv.ParseInt(os.Getenv("REDIS_TTL_KEY"), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("enviroment REDIS_TTL_KEY not a number")
		}
		keyTtl = time.Duration(keyTtlInt)

	}
	return &RedisConf{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB_Index: dbIndex,
		KeyTtl:   keyTtl,
	}, nil
}

func loadEnv(pathToEnv string) error {
	if err := godotenv.Load(pathToEnv); err != nil {
		return fmt.Errorf("failed to load %s: %w", pathToEnv, err)
	}
	return nil
}
