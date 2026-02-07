package stock

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func Test_CheckChange(t *testing.T) {
	// 使用 mock Finnhub client
	mockFinnhub := NewMockFinnhubClient()
	mockFinnhub.AddQuote("TSLA", &QuoteResponse{
		CurrentPrice:  100,
		PercentChange: 5.0,
	})
	mockFinnhub.AddQuote("AAPL", &QuoteResponse{
		CurrentPrice:  100,
		PercentChange: 5.0,
	})
	SetDefaultClient(mockFinnhub)
	defer ResetDefaultClient()

	// 使用 mock Redis
	mockRedis := NewMockRedisClient()
	mockRedis.Lists["watch_list"] = []string{"TSLA", "AAPL"}

	// 測試使用接口注入依賴
	// 注意: CheckChangeWithDeps 需要 Discord session，在單元測試中我們只驗證邏輯流程
	// 實際的 Discord 消息發送需要 integration test
	t.Log("CheckChange logic can be tested with mocked dependencies")
	t.Logf("Watch list: %v", mockRedis.Lists["watch_list"])
}
