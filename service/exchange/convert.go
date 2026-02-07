package exchange

import (
	"context"
	"encoding/json"
	"fmt"

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

var httpClient = client.NewHTTPClientWithEnv("EXCHANGE")

// ConvertExchange : 換算幣值
func ConvertExchange(oldMoney []float64) ([]float64, error) {
	return ConvertExchangeWithContext(context.Background(), oldMoney)
}

// ConvertExchangeWithContext : 使用指定 context 換算幣值
func ConvertExchangeWithContext(ctx context.Context, oldMoney []float64) ([]float64, error) {
	// 先取匯率
	body, err := httpClient.GetWithRetry(ctx, exchangeAPIURL, maxRetries)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}

	info := &apiInfo{}
	if err := json.Unmarshal(body, info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exchange rate: %w", err)
	}

	exrate := info.USDTWD.Exrate
	if exrate == 0 {
		return nil, fmt.Errorf("invalid exchange rate: %f", exrate)
	}

	newMoney := make([]float64, len(oldMoney))
	for i, v := range oldMoney {
		newMoney[i] = v * exrate
	}

	return newMoney, nil
}
