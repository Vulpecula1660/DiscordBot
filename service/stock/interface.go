package stock

import (
	"context"
)

// FinnhubClient Finnhub API 客戶端接口
type FinnhubClient interface {
	GetQuote(ctx context.Context, symbol string) (*QuoteResponse, error)
}

// QuoteResponse 報價回應
type QuoteResponse struct {
	CurrentPrice  float32
	PercentChange float32
	Change        float32
	HighPrice     float32
	LowPrice      float32
	OpenPrice     float32
	PreviousClose float32
}
