package stock

import (
	"context"
	"time"

	stockdao "discordBot/model/dao/stock"
	"discordBot/model/dto"
)

// MockFinnhubClient FinnhubClient 的 mock 實現
type MockFinnhubClient struct {
	Quotes map[string]*QuoteResponse
	Err    error
}

// GetQuote 實現 FinnhubClient 接口
func (m *MockFinnhubClient) GetQuote(ctx context.Context, symbol string) (*QuoteResponse, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if quote, ok := m.Quotes[symbol]; ok {
		return quote, nil
	}

	// 默認返回空報價（表示搜尋失敗）
	return &QuoteResponse{}, nil
}

// NewMockFinnhubClient 創建一個新的 mock client
func NewMockFinnhubClient() *MockFinnhubClient {
	return &MockFinnhubClient{
		Quotes: make(map[string]*QuoteResponse),
	}
}

// AddQuote 添加報價到 mock
func (m *MockFinnhubClient) AddQuote(symbol string, quote *QuoteResponse) {
	m.Quotes[symbol] = quote
}

// MockRedisClient Redis Client 的 mock 實現
type MockRedisClient struct {
	Data  map[string]string
	Lists map[string][]string
	Err   error
}

// NewMockRedisClient 創建新的 mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		Data:  make(map[string]string),
		Lists: make(map[string][]string),
	}
}

// Get 實現 Client 接口
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.Data[key], nil
}

// Set 實現 Client 接口
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if m.Err != nil {
		return m.Err
	}
	m.Data[key] = value.(string)
	return nil
}

// LRange 實現 Client 接口
func (m *MockRedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	list := m.Lists[key]
	if list == nil {
		return []string{}, nil
	}
	return list, nil
}

// MockStockRepository Stock Repository 的 mock 實現
type MockStockRepository struct {
	Stocks []*dto.Stock
	Err    error
}

// Get 實現 Repository 接口
func (m *MockStockRepository) Get(ctx context.Context, input *stockdao.GetInput) ([]*dto.Stock, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Stocks, nil
}
