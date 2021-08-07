package redis

import (
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type (
	// 連線池條件
	redisPoolCond struct {
		name string
		db   int
	}
)

var (
	// Redis連線物件
	pool = make(map[redisPoolCond]*redis.Client)

	//同步鎖
	mu *sync.Mutex
)

func init() {
	pool = make(map[redisPoolCond]*redis.Client)
	mu = &sync.Mutex{}
}

// GetConn : 取得redis連線
// 若超時無法取得連線，會回傳error
func GetConn(name string, db int) (ret *redis.Client) {
	mu.Lock()
	defer mu.Unlock()

	cond := redisPoolCond{name, db}

	if conn, ok := pool[cond]; ok {
		return conn
	}

	pool[cond] = CreateConn(name, db)
	return pool[cond]

}

// CreateConn : 建立redis連線
func CreateConn(name string, db int) (ret *redis.Client) {

	redisHost := os.Getenv("RedisHost")
	redisPassword := os.Getenv("RedisPassword")

	return redis.NewClient(
		&redis.Options{
			ReadTimeout:        time.Second * 3,
			WriteTimeout:       time.Second * 3,
			DB:                 db,
			Addr:               redisHost,
			Password:           redisPassword,
			MaxRetries:         15,
			DialTimeout:        time.Second * 10,
			PoolSize:           50,
			PoolTimeout:        time.Second * 5,
			IdleTimeout:        time.Minute * 5,
			IdleCheckFrequency: time.Minute,
		},
	)
}
