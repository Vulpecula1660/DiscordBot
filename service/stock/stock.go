package stock

import (
	"context"
	"fmt"
	"os"
	"strings"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type QuoteRes struct {
	// Open price of the day
	O *float32 `json:"o,omitempty"`
	// High price of the day
	H *float32 `json:"h,omitempty"`
	// Low price of the day
	L *float32 `json:"l,omitempty"`
	// Current price
	C *float32 `json:"c,omitempty"`
	// Previous close price
	Pc *float32 `json:"pc,omitempty"`
	// Change
	D *float32 `json:"d,omitempty"`
	// Percent change
	Dp *float32 `json:"dp,omitempty"`
}

func Quote(message string) (string, error) {

	strSlice := strings.Split(message, "$+")

	symbol := strings.ToUpper(strSlice[1])

	key := os.Getenv("APIKey")

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", key)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	res, _, err := finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return "", err
	}

	if *res.C == 0 {
		return "", fmt.Errorf("搜尋失敗")
	}

	resStr1 := fmt.Sprintf("Finnhub 查詢標的為:%s 目前價格為:%v 今天漲跌幅:%v%s", symbol, *res.C, *res.Dp, "%")

	return resStr1, err
}

func GetChange(stock string) (float32, error) {
	key := os.Getenv("APIKey")

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", key)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	res, _, err := finnhubClient.Quote(context.Background()).Symbol(stock).Execute()

	if err != nil {
		return 0, err
	}

	if *res.C == 0 {
		return 0, fmt.Errorf("搜尋失敗")
	}

	return *res.Dp, err
}
