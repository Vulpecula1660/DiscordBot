package crypto

import (
	"context"
	"encoding/json"
	"fmt"

	"discordBot/service/client"
)

type (
	apiInfo struct {
		Ethereum usd `json:"ethereum"`
	}

	usd struct {
		Price float64 `json:"usd"`
	}
)

const (
	coingeckoAPIURL = "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
	maxRetries      = 3
)

var httpClient = client.NewHTTPClientWithEnv("CRYPTO")

// GetPrice 獲取ETH當前價格
func GetPrice() (float64, error) {
	return GetPriceWithContext(context.Background())
}

// GetPriceWithContext 使用指定 context 獲取ETH當前價格
func GetPriceWithContext(ctx context.Context) (float64, error) {
	body, err := httpClient.GetWithRetry(ctx, coingeckoAPIURL, maxRetries)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch crypto price: %w", err)
	}

	info := &apiInfo{
		Ethereum: usd{},
	}

	if err := json.Unmarshal(body, info); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return info.Ethereum.Price, nil
}
