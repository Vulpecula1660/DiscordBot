package exchange

import (
	"context"

	"discordBot/service/client"
)

type (
	apiInfo struct {
		USDTWD usdtwd `json:"USDTWD"`
	}

	usdtwd struct {
		Exrate float64 `json:"Exrate"`
	}
)

const (
	exchangeAPIURL = "https://tw.rter.info/capi.php"
	maxRetries     = 3
)

var (
	httpClient   = client.NewHTTPClientWithEnv("EXCHANGE")
	rateProvider ExchangeRateProvider
)

func init() {
	rateProvider = NewRateProvider(httpClientAdapter{client: httpClient})
}

// httpClientAdapter 適配器將 HTTPClient 轉換為 HTTPClient 接口
type httpClientAdapter struct {
	client *client.HTTPClient
}

func (a httpClientAdapter) GetWithRetry(ctx context.Context, url string, maxRetries int) ([]byte, error) {
	return a.client.GetWithRetry(ctx, url, maxRetries)
}

// SetRateProvider 設置匯率提供者（用於測試）
func SetRateProvider(provider ExchangeRateProvider) {
	rateProvider = provider
}

// ResetRateProvider 重置匯率提供者為默認值
func ResetRateProvider() {
	rateProvider = NewRateProvider(httpClientAdapter{client: httpClient})
}

// ConvertExchange : 換算幣值
func ConvertExchange(oldMoney []float64) ([]float64, error) {
	return ConvertExchangeWithContext(context.Background(), oldMoney)
}

// ConvertExchangeWithContext : 使用指定 context 換算幣值
func ConvertExchangeWithContext(ctx context.Context, oldMoney []float64) ([]float64, error) {
	exrate, err := rateProvider.GetRate(ctx)
	if err != nil {
		return nil, err
	}

	newMoney := make([]float64, len(oldMoney))
	for i, v := range oldMoney {
		newMoney[i] = v * exrate
	}

	return newMoney, nil
}
