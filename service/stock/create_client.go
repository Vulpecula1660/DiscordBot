package stock

import (
	"context"
	"net/http"
	"os"
	"sync"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var (
	// Finnhub 連線物件
	pool map[string]FinnhubClient

	//同步鎖
	mu *sync.Mutex

	// defaultClientFactory 用於建立 Finnhub client，可在測試中替換
	defaultClientFactory = createRealClient
)

func init() {
	pool = make(map[string]FinnhubClient)
	mu = &sync.Mutex{}
}

// finnhubClientAdapter 包裝 Finnhub DefaultApiService 以實現 FinnhubClient 介面
type finnhubClientAdapter struct {
	client *finnhub.DefaultApiService
}

// Quote 實現 FinnhubClient 介面
func (a *finnhubClientAdapter) Quote(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
	return a.client.Quote(ctx).Symbol(symbol).Execute()
}

// GetConn : 取得 Finnhub 連線
// 若超時無法取得連線，會回傳error
func GetConn(name string) FinnhubClient {
	mu.Lock()
	defer mu.Unlock()

	if conn, ok := pool[name]; ok {
		return conn
	}

	pool[name] = defaultClientFactory(name)
	return pool[name]
}

// createRealClient : 建立真實的 Finnhub 連線
func createRealClient(name string) FinnhubClient {
	key := os.Getenv("APIKey")

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", key)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	return &finnhubClientAdapter{client: finnhubClient}
}

// SetClientFactory 設定 client factory，用於測試時注入 mock
func SetClientFactory(factory func(string) FinnhubClient) {
	mu.Lock()
	defer mu.Unlock()
	defaultClientFactory = factory
	// 清空現有連線池，強制使用新的 factory
	pool = make(map[string]FinnhubClient)
}

// ResetClientFactory 重置為預設的 client factory
func ResetClientFactory() {
	mu.Lock()
	defer mu.Unlock()
	defaultClientFactory = createRealClient
	pool = make(map[string]FinnhubClient)
}
