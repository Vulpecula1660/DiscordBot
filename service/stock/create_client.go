package stock

import (
	"os"
	"sync"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var (
	// Finnhub 連線物件
	pool map[string]*finnhub.DefaultApiService

	//同步鎖
	mu *sync.Mutex
)

func init() {
	pool = make(map[string]*finnhub.DefaultApiService)
	mu = &sync.Mutex{}
}

// GetConn : 取得 Finnhub 連線
// 若超時無法取得連線，會回傳error
func GetConn(name string) (ret *finnhub.DefaultApiService) {
	mu.Lock()
	defer mu.Unlock()

	if conn, ok := pool[name]; ok {
		return conn
	}

	pool[name] = createConn(name)
	return pool[name]
}

// CreateConn : 建立 Finnhub 連線
func createConn(name string) (ret *finnhub.DefaultApiService) {

	key := os.Getenv("APIKey")

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", key)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	return finnhubClient
}
