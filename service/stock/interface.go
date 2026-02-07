package stock

import (
	"context"

	"discordBot/model/dto"
)

// Quoter 股票報價接口
type Quoter interface {
	Quote(ctx context.Context, symbol string) (*QuoteResult, error)
	GetChange(ctx context.Context, symbol string) (float32, error)
}

// QuoteResult 報價結果
type QuoteResult struct {
	Symbol string
	Price  float64
	Change float32
}

// Calculator 股票計算器接口
type Calculator interface {
	Calculate(ctx context.Context, input *CalculateInput) (*CalculateResult, error)
}

// CalculateResult 計算結果
type CalculateResult struct {
	Value  float64
	Profit float64
}

// Repository 股票數據倉庫接口
type Repository interface {
	Get(ctx context.Context, userID string, symbol string) ([]*dto.Stock, error)
	Save(ctx context.Context, stock *dto.Stock) error
}
