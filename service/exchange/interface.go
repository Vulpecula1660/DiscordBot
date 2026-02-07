package exchange

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExchangeRateProvider 匯率提供者接口
type ExchangeRateProvider interface {
	GetRate(ctx context.Context) (float64, error)
}

// HTTPClient HTTP 客戶端接口
type HTTPClient interface {
	GetWithRetry(ctx context.Context, url string, maxRetries int) ([]byte, error)
}

// defaultRateProvider 默認匯率提供者實現
type defaultRateProvider struct {
	client HTTPClient
}

// NewRateProvider 創建匯率提供者
func NewRateProvider(client HTTPClient) ExchangeRateProvider {
	return &defaultRateProvider{client: client}
}

// GetRate 獲取匯率
func (p *defaultRateProvider) GetRate(ctx context.Context) (float64, error) {
	body, err := p.client.GetWithRetry(ctx, exchangeAPIURL, maxRetries)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}

	info := &apiInfo{}
	if err := json.Unmarshal(body, info); err != nil {
		return 0, fmt.Errorf("failed to unmarshal exchange rate: %w", err)
	}

	exrate := info.USDTWD.Exrate
	if exrate == 0 {
		return 0, fmt.Errorf("invalid exchange rate: %f", exrate)
	}

	return exrate, nil
}
