package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"discordBot/pkg/logger"
)

var (
	// Redis 連線物件
	client *redis.Client
	mu     sync.Mutex
)

// getClient : 取得 Redis 連線
func getClient() (*redis.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		return client, nil
	}

	conn, err := createConn()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis connection: %w", err)
	}

	client = conn
	return client, nil
}

// createConn : 建立 Redis 連線
func createConn() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_URL")
	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_URL environment variable is not set")
	}

	opt, err := redis.ParseURL(redisHost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	// 配置連接池
	opt.PoolSize = getEnvInt("REDIS_POOL_SIZE", 10)
	opt.MinIdleConns = getEnvInt("REDIS_MIN_IDLE_CONNS", 5)
	opt.MaxRetries = getEnvInt("REDIS_MAX_RETRIES", 3)
	opt.ReadTimeout = getEnvDuration("REDIS_READ_TIMEOUT", 10*time.Second)
	opt.WriteTimeout = getEnvDuration("REDIS_WRITE_TIMEOUT", 10*time.Second)
	opt.PoolTimeout = getEnvDuration("REDIS_POOL_TIMEOUT", 30*time.Second)

	c := redis.NewClient(opt)

	// 驗證連線
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return c, nil
}

// Close : 關閉 Redis 連線
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		if err := client.Close(); err != nil {
			logger.Error("關閉Redis連線失敗", "error", err)
			return fmt.Errorf("failed to close redis connection: %w", err)
		}
		client = nil
	}
	return nil
}

// getEnvInt : 從環境變量獲取整數值，帶有默認值
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvDuration : 從環境變量獲取時間值，帶有默認值
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultVal
}
