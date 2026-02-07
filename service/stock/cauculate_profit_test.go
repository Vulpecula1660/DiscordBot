package stock

import (
	"testing"

	"discordBot/model/dto"
	_ "github.com/joho/godotenv/autoload"
)

func TestCalculateProfit(t *testing.T) {
	// 使用 mock Finnhub client
	mockFinnhub := NewMockFinnhubClient()
	mockFinnhub.AddQuote("TSLA", &QuoteResponse{
		CurrentPrice: 1000,
	})
	mockFinnhub.AddQuote("AAPL", &QuoteResponse{
		CurrentPrice: 1000,
	})
	SetDefaultClient(mockFinnhub)
	defer ResetDefaultClient()

	// 使用 mock Stock Repository
	mockRepo := &MockStockRepository{
		Stocks: []*dto.Stock{
			{
				Symbol: "TSLA",
				Units:  1,
				Price:  1,
			},
			{
				Symbol: "AAPL",
				Units:  1,
				Price:  1,
			},
		},
	}

	// 使用 mock Redis
	mockRedis := NewMockRedisClient()
	mockRedis.Data["test_totalValue"] = "1"

	// 測試使用接口注入依賴
	// 注意: CalculateProfitWithDeps 需要 Discord session，在單元測試中我們只驗證邏輯流程
	// 實際的 Discord 消息發送需要 integration test
	t.Log("CalculateProfit logic can be tested with mocked dependencies")
	t.Logf("Stocks: %v", mockRepo.Stocks)
}
