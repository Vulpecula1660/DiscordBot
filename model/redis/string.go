package redis

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

const (
	// redisDB : 放在 redis db0
	redisDB int = 0
)

// Set : 資料寫入Redis中
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {

	return CreateConn().Set(
		ctx,
		key,
		value,
		expiration,
	).Err()
}

// Get : 從 Redis 取得資料
func Get(ctx context.Context, key string) (string, error) {

	data, err := CreateConn().Get(
		ctx,
		key,
	).Result()

	// 找不到Key
	if err == redis.Nil || data == "" {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return data, nil
}

// Del :
func Del(ctx context.Context, key string) error {

	return CreateConn().
		Del(
			ctx,
			key,
		).Err()
}
