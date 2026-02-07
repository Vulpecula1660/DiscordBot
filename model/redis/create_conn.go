package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// Redis連線物件
	pool map[string]*redis.Client

	// 讀寫鎖，減少讀取路徑的競爭
	mu sync.RWMutex
)

func init() {
	pool = make(map[string]*redis.Client)
}

// GetConn : 取得redis連線
func GetConn(name string) (*redis.Client, error) {
	// 快速路徑：讀鎖取得已存在的連線
	mu.RLock()
	conn, ok := pool[name]
	mu.RUnlock()

	if ok {
		// 在鎖外驗證連線狀態，避免阻塞其他 goroutine
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := conn.Ping(ctx).Err(); err == nil {
			return conn, nil
		}

		// 連線失敗，需要重建
		mu.Lock()
		// Double-check：其他 goroutine 可能已經重建
		if currentConn, exists := pool[name]; exists && currentConn == conn {
			conn.Close()
			delete(pool, name)
		}
		mu.Unlock()
	}

	// 慢速路徑：寫鎖建立新連線
	mu.Lock()
	defer mu.Unlock()

	// Double-check：其他 goroutine 可能已經建立
	if conn, ok := pool[name]; ok {
		return conn, nil
	}

	conn, err := createConn(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis connection: %w", err)
	}

	pool[name] = conn
	return conn, nil
}

// createConn : 建立redis連線
func createConn(name string) (*redis.Client, error) {
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

	client := redis.NewClient(opt)

	// 驗證連線
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

// CloseAll : 關閉所有Redis連線
func CloseAll() error {
	mu.Lock()
	defer mu.Unlock()

	var errs []error
	for name, conn := range pool {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close redis connection %s: %w", name, err))
		}
		delete(pool, name)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing redis connections: %v", errs)
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
