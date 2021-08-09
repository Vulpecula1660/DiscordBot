package redis

import (
	"os"

	"github.com/go-redis/redis/v8"
)

// CreateConn : 建立redis連線
func CreateConn() (ret *redis.Client) {

	redisHost := os.Getenv("REDIS_URL")

	opt, _ := redis.ParseURL(redisHost)

	return redis.NewClient(
		opt,
	)
}
