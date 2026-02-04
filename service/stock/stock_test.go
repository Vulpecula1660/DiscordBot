package stock

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

// MockFinnhubClient 實現 FinnhubClient 介面的 mock
type MockFinnhubClient struct {
	QuoteFunc func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error)
}

func (m *MockFinnhubClient) Quote(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
	if m.QuoteFunc != nil {
		return m.QuoteFunc(ctx, symbol)
	}
	return finnhub.Quote{}, nil, fmt.Errorf("QuoteFunc not implemented")
}

// newMockQuote 建立一個 mock Quote 回應
func newMockQuote(currentPrice, percentChange float32) finnhub.Quote {
	c := currentPrice
	dp := percentChange
	return finnhub.Quote{
		C:  &c,
		Dp: &dp,
	}
}

// setupMockClient 設定 mock client 並返回清理函數
func setupMockClient(mock *MockFinnhubClient) func() {
	SetClientFactory(func(name string) FinnhubClient {
		return mock
	})
	return func() {
		ResetClientFactory()
	}
}

func Test_Quote(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		message   string
		mockQuote func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error)
		want      string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "TSLA Success",
			message: "$+TSLA",
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return newMockQuote(250.50, 3.25), nil, nil
			},
			want:    "Finnhub 查詢標的為:+TSLA 目前價格為:250.5 今天漲跌幅:3.25%",
			wantErr: false,
		},
		{
			name:    "Wrong Message Format",
			message: "aaaaa",
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return finnhub.Quote{}, nil, nil
			},
			want:    "",
			wantErr: true,
			errMsg:  "參數錯誤",
		},
		{
			name:    "Symbol Not Found",
			message: "$+INVALID",
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				// 返回價格為 0 表示找不到
				return newMockQuote(0, 0), nil, nil
			},
			want:    "",
			wantErr: true,
			errMsg:  "搜尋失敗",
		},
		{
			name:    "API Error",
			message: "$+TSLA",
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return finnhub.Quote{}, nil, fmt.Errorf("API connection error")
			},
			want:    "",
			wantErr: true,
			errMsg:  "API connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設定 mock
			cleanup := setupMockClient(&MockFinnhubClient{
				QuoteFunc: tt.mockQuote,
			})
			defer cleanup()

			got, err := Quote(ctx, tt.message)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Quote() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Quote() error = %v, wantErr %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Quote() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetChange(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		stock     string
		mockQuote func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error)
		want      float32
		wantErr   bool
		errMsg    string
	}{
		{
			name:  "TSLA Success",
			stock: "TSLA",
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return newMockQuote(250.50, 3.25), nil, nil
			},
			want:    3.25,
			wantErr: false,
		},
		{
			name:  "Symbol Not Found",
			stock: "INVALID",
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return newMockQuote(0, 0), nil, nil
			},
			want:    0,
			wantErr: true,
			errMsg:  "搜尋失敗",
		},
		{
			name:  "API Error",
			stock: "TSLA",
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return finnhub.Quote{}, nil, fmt.Errorf("API connection error")
			},
			want:    0,
			wantErr: true,
			errMsg:  "API connection error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupMockClient(&MockFinnhubClient{
				QuoteFunc: tt.mockQuote,
			})
			defer cleanup()

			got, err := GetChange(ctx, tt.stock)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetChange() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("GetChange() error = %v, wantErr %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("GetChange() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("GetChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Calculate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		input      *CalculateInput
		mockQuote  func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error)
		wantValue  float64
		wantProfit float64
		wantErr    bool
		errMsg     string
	}{
		{
			name: "TSLA Success",
			input: &CalculateInput{
				Symbol: "TSLA",
				Units:  10,
				Price:  200,
			},
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return newMockQuote(250, 5.0), nil, nil
			},
			wantValue:  2500, // 10 * 250
			wantProfit: 500,  // 2500 - (10 * 200)
			wantErr:    false,
		},
		{
			name: "Symbol Not Found",
			input: &CalculateInput{
				Symbol: "INVALID",
				Units:  1,
				Price:  1,
			},
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return newMockQuote(0, 0), nil, nil
			},
			wantValue:  0,
			wantProfit: 0,
			wantErr:    true,
			errMsg:     "搜尋失敗",
		},
		{
			name: "API Error",
			input: &CalculateInput{
				Symbol: "TSLA",
				Units:  1,
				Price:  1,
			},
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return finnhub.Quote{}, nil, fmt.Errorf("API connection error")
			},
			wantValue:  0,
			wantProfit: 0,
			wantErr:    true,
			errMsg:     "API connection error",
		},
		{
			name: "Negative Profit (Loss)",
			input: &CalculateInput{
				Symbol: "TSLA",
				Units:  10,
				Price:  300,
			},
			mockQuote: func(ctx context.Context, symbol string) (finnhub.Quote, *http.Response, error) {
				return newMockQuote(250, -5.0), nil, nil
			},
			wantValue:  2500, // 10 * 250
			wantProfit: -500, // 2500 - (10 * 300)
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupMockClient(&MockFinnhubClient{
				QuoteFunc: tt.mockQuote,
			})
			defer cleanup()

			gotValue, gotProfit, err := Calculate(ctx, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Calculate() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Calculate() error = %v, wantErr %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Calculate() unexpected error = %v", err)
				return
			}

			if gotValue != tt.wantValue {
				t.Errorf("Calculate() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotProfit != tt.wantProfit {
				t.Errorf("Calculate() gotProfit = %v, want %v", gotProfit, tt.wantProfit)
			}
		})
	}
}
