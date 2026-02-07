package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Set : 資料寫入Redis中
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	conn, err := getClient()
	if err != nil {
		return err
	}

	return conn.Set(ctx, key, value, expiration).Err()
}

// Get : 從 Redis 取得資料
func Get(ctx context.Context, key string) (string, error) {
	conn, err := getClient()
	if err != nil {
		return "", err
	}

	data, err := conn.Get(ctx, key).Result()

	// 找不到Key
	if err == redis.Nil || data == "" {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return data, nil
}

// Del : 從 Redis 刪除資料
func Del(ctx context.Context, key string) error {
	conn, err := getClient()
	if err != nil {
		return err
	}

	return conn.Del(ctx, key).Err()
}
