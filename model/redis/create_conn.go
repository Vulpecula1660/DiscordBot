package redis

import (
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	// Redis連線物件
	pool map[string]*redis.Client

	//同步鎖
	mu *sync.Mutex
)

func init() {
	pool = make(map[string]*redis.Client)
	mu = &sync.Mutex{}
}

// GetConn : 取得redis連線
// 若超時無法取得連線，會回傳error
func GetConn(name string) (ret *redis.Client) {
	mu.Lock()
	defer mu.Unlock()

	if conn, ok := pool[name]; ok {
		return conn
	}

	pool[name] = createConn(name)
	return pool[name]
}

// CreateConn : 建立redis連線
func createConn(name string) (ret *redis.Client) {

	redisHost := os.Getenv("REDIS_URL")

	opt, _ := redis.ParseURL(redisHost)

	return redis.NewClient(
		opt,
	)
}
