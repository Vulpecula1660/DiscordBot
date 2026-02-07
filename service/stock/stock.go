package stock

import (
	"context"
	"fmt"
	"strings"

	"discordBot/pkg/logger"
)

// Quote : 查詢標的
func Quote(ctx context.Context, message string) (string, error) {
	// example : $+TSLA
	strSlice := strings.Split(message, "$")

	if len(strSlice) != 2 {
		logger.Warn("股票查詢參數錯誤", "message", message)
		return "", fmt.Errorf("參數錯誤")
	}

	symbol := strings.ToUpper(strSlice[1])
	logger.Info("查詢股票價格", "symbol", symbol)

	finnhubClient := GetConn("finnhub")

	res, _, err := finnhubClient.Quote(ctx).Symbol(symbol).Execute()
	if err != nil {
		logger.Error("查詢股票價格失敗", "symbol", symbol, "error", err)
		return "", err
	}

	// GetC : Get Current price
	if res.GetC() == 0 {
		logger.Warn("股票查詢結果為空", "symbol", symbol)
		return "", fmt.Errorf("搜尋失敗")
	}

	resStr1 := fmt.Sprintf("Finnhub 查詢標的為:%s 目前價格為:%v 今天漲跌幅:%v%s", symbol, *res.C, *res.Dp, "%")
	logger.Info("股票查詢成功", "symbol", symbol, "price", *res.C, "change", *res.Dp)

	return resStr1, nil
}

// GetChange : 取得漲跌幅
func GetChange(ctx context.Context, stock string) (float32, error) {
	finnhubClient := GetConn("finnhub")

	symbol := strings.ToUpper(stock)

	res, _, err := finnhubClient.Quote(ctx).Symbol(symbol).Execute()

	if err != nil {
		return 0, err
	}

	// GetC : Get Current price
	if res.GetC() == 0 {
		return 0, fmt.Errorf("搜尋失敗")
	}

	// GetDp : Get Percent change
	return res.GetDp(), nil
}

type CalculateInput struct {
	Symbol string
	Units  float64
	Price  float64
}

// Calculate : 計算成本, 損益
func Calculate(ctx context.Context, input *CalculateInput) (value, profit float64, err error) {

	finnhubClient := GetConn("finnhub")

	res, _, err := finnhubClient.Quote(ctx).Symbol(input.Symbol).Execute()
	if err != nil {
		return 0, 0, err
	}

	// Current price
	c := res.GetC()

	if c == 0 {
		return 0, 0, fmt.Errorf("搜尋失敗")
	}

	// 成本
	cost := input.Units * input.Price

	// 市場價值
	value = input.Units * float64(c)

	// 損益
	profit = value - cost

	return value, profit, nil
}
