package exchange

import "context"

// MockRateProvider ExchangeRateProvider 的 mock 實現
type MockRateProvider struct {
	Rate float64
	Err  error
}

// GetRate 實現 ExchangeRateProvider 接口
func (m *MockRateProvider) GetRate(ctx context.Context) (float64, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return m.Rate, nil
}
