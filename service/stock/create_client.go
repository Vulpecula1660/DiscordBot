package stock

import (
	"context"
	"os"
	"sync"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var (
	// Finnhub 連線物件
	pool map[string]*finnhub.DefaultApiService

	// clientPool 儲存 FinnhubClient 接口實例
	clientPool map[string]FinnhubClient

	// 同步鎖
	mu *sync.Mutex

	// defaultClient 默認客戶端（用於測試替換）
	defaultClient FinnhubClient
)

func init() {
	pool = make(map[string]*finnhub.DefaultApiService)
	clientPool = make(map[string]FinnhubClient)
	mu = &sync.Mutex{}
}

// finnhubClientWrapper 包裝 finnhub client 以實現 FinnhubClient 接口
type finnhubClientWrapper struct {
	client *finnhub.DefaultApiService
}

// GetQuote 實現 FinnhubClient 接口
func (w *finnhubClientWrapper) GetQuote(ctx context.Context, symbol string) (*QuoteResponse, error) {
	res, _, err := w.client.Quote(ctx).Symbol(symbol).Execute()
	if err != nil {
		return nil, err
	}

	return &QuoteResponse{
		CurrentPrice:  res.GetC(),
		PercentChange: res.GetDp(),
		Change:        res.GetD(),
		HighPrice:     res.GetH(),
		LowPrice:      res.GetL(),
		OpenPrice:     res.GetO(),
		PreviousClose: res.GetPc(),
	}, nil
}

// GetConn : 取得 Finnhub 連線（向後兼容）
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

// GetClient : 取得 FinnhubClient 接口實例
func GetClient(name string) FinnhubClient {
	mu.Lock()
	defer mu.Unlock()

	if defaultClient != nil {
		return defaultClient
	}

	if client, ok := clientPool[name]; ok {
		return client
	}

	clientPool[name] = &finnhubClientWrapper{client: createConn(name)}
	return clientPool[name]
}

// SetDefaultClient : 設置默認客戶端（用於測試）
func SetDefaultClient(client FinnhubClient) {
	mu.Lock()
	defer mu.Unlock()
	defaultClient = client
}

// ResetDefaultClient : 重置默認客戶端
func ResetDefaultClient() {
	mu.Lock()
	defer mu.Unlock()
	defaultClient = nil
}

// CreateConn : 建立 Finnhub 連線
func createConn(name string) (ret *finnhub.DefaultApiService) {
	key := os.Getenv("APIKey")

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", key)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	return finnhubClient
}
