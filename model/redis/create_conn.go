package redis

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"discordBot/pkg/config"
)

var (
	// Redis連線物件
	pool map[string]*redis.Client

	// 同步鎖
	mu *sync.Mutex
)

func init() {
	pool = make(map[string]*redis.Client)
	mu = &sync.Mutex{}
}

// GetConn : 取得redis連線
func GetConn(name string) (*redis.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if conn, ok := pool[name]; ok {
		// 驗證連線狀態
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := conn.Ping(ctx).Err(); err != nil {
			// 連線失敗，重新建立
			conn.Close()
			delete(pool, name)
		} else {
			return conn, nil
		}
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
	opt.PoolSize = config.GetEnvInt("REDIS_POOL_SIZE", 10)
	opt.MinIdleConns = config.GetEnvInt("REDIS_MIN_IDLE_CONNS", 5)
	opt.MaxRetries = config.GetEnvInt("REDIS_MAX_RETRIES", 3)
	opt.ReadTimeout = config.GetEnvDuration("REDIS_READ_TIMEOUT", 10*time.Second)
	opt.WriteTimeout = config.GetEnvDuration("REDIS_WRITE_TIMEOUT", 10*time.Second)
	opt.PoolTimeout = config.GetEnvDuration("REDIS_POOL_TIMEOUT", 30*time.Second)

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
