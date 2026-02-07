package stock

import (
	"context"
	"errors"
	"testing"
)

func Test_Quote(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		message   string
		mockSetup func(*MockFinnhubClient)
		want      string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "successful query - TSLA",
			message: "$+TSLA",
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("+TSLA", &QuoteResponse{
					CurrentPrice:  1137.06,
					PercentChange: 3.7104,
				})
			},
			want:    "Finnhub 查詢標的為:+TSLA 目前價格為:1137.06 今天漲跌幅:3.7104%",
			wantErr: false,
		},
		{
			name:      "invalid message format - no dollar sign",
			message:   "aaaaa",
			mockSetup: func(m *MockFinnhubClient) {},
			want:      "",
			wantErr:   true,
			errMsg:    "參數錯誤",
		},
		{
			name:    "symbol with plus prefix",
			message: "$+TSLA",
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("+TSLA", &QuoteResponse{
					CurrentPrice:  200.0,
					PercentChange: 2.5,
				})
			},
			want:    "Finnhub 查詢標的為:+TSLA 目前價格為:200 今天漲跌幅:2.5%",
			wantErr: false,
		},
		{
			name:    "symbol not found - zero price",
			message: "$+INVALID",
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("INVALID", &QuoteResponse{
					CurrentPrice: 0,
				})
			},
			want:    "",
			wantErr: true,
			errMsg:  "搜尋失敗",
		},
		{
			name:    "API error",
			message: "$+AAPL",
			mockSetup: func(m *MockFinnhubClient) {
				m.Err = errors.New("API rate limit exceeded")
			},
			want:    "",
			wantErr: true,
			errMsg:  "API rate limit exceeded",
		},
		{
			name:    "case insensitive symbol with plus",
			message: "$+tsla",
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("+TSLA", &QuoteResponse{
					CurrentPrice:  100.0,
					PercentChange: 1.5,
				})
			},
			want:    "Finnhub 查詢標的為:+TSLA 目前價格為:100 今天漲跌幅:1.5%",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設置 mock
			mock := NewMockFinnhubClient()
			tt.mockSetup(mock)
			SetDefaultClient(mock)
			defer ResetDefaultClient()

			got, err := Quote(ctx, tt.message)

			// 驗證錯誤
			if tt.wantErr {
				if err == nil {
					t.Errorf("Quote() error = nil, wantErr = true")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Quote() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Quote() unexpected error = %v", err)
				return
			}

			// 驗證結果
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
		mockSetup func(*MockFinnhubClient)
		want      float32
		wantErr   bool
		errMsg    string
	}{
		{
			name:  "successful query - TSLA",
			stock: "TSLA",
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("TSLA", &QuoteResponse{
					CurrentPrice:  100.0,
					PercentChange: 3.7104,
				})
			},
			want:    3.7104,
			wantErr: false,
		},
		{
			name:  "symbol not found - zero price",
			stock: "INVALID",
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("INVALID", &QuoteResponse{
					CurrentPrice: 0,
				})
			},
			want:    0,
			wantErr: true,
			errMsg:  "搜尋失敗",
		},
		{
			name:  "API error",
			stock: "AAPL",
			mockSetup: func(m *MockFinnhubClient) {
				m.Err = errors.New("connection timeout")
			},
			want:    0,
			wantErr: true,
			errMsg:  "connection timeout",
		},
		{
			name:  "case insensitive",
			stock: "tsla",
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("TSLA", &QuoteResponse{
					CurrentPrice:  100.0,
					PercentChange: 5.0,
				})
			},
			want:    5.0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設置 mock
			mock := NewMockFinnhubClient()
			tt.mockSetup(mock)
			SetDefaultClient(mock)
			defer ResetDefaultClient()

			got, err := GetChange(ctx, tt.stock)

			// 驗證錯誤
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetChange() error = nil, wantErr = true")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("GetChange() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("GetChange() unexpected error = %v", err)
				return
			}

			// 驗證結果
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
		mockSetup  func(*MockFinnhubClient)
		wantValue  float64
		wantProfit float64
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful calculation",
			input: &CalculateInput{
				Symbol: "TSLA",
				Units:  10,
				Price:  100,
			},
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("TSLA", &QuoteResponse{
					CurrentPrice: 113.706,
				})
			},
			wantValue:  1137.06,
			wantProfit: 137.06,
			wantErr:    false,
		},
		{
			name: "loss calculation",
			input: &CalculateInput{
				Symbol: "AAPL",
				Units:  5,
				Price:  200,
			},
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("AAPL", &QuoteResponse{
					CurrentPrice: 180,
				})
			},
			wantValue:  900,
			wantProfit: -100,
			wantErr:    false,
		},
		{
			name: "symbol not found",
			input: &CalculateInput{
				Symbol: "INVALID",
				Units:  1,
				Price:  1,
			},
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("INVALID", &QuoteResponse{
					CurrentPrice: 0,
				})
			},
			wantValue:  0,
			wantProfit: 0,
			wantErr:    true,
			errMsg:     "搜尋失敗",
		},
		{
			name: "API error",
			input: &CalculateInput{
				Symbol: "ERROR",
				Units:  1,
				Price:  1,
			},
			mockSetup: func(m *MockFinnhubClient) {
				m.Err = errors.New("service unavailable")
			},
			wantValue:  0,
			wantProfit: 0,
			wantErr:    true,
			errMsg:     "service unavailable",
		},
		{
			name: "zero units",
			input: &CalculateInput{
				Symbol: "TSLA",
				Units:  0,
				Price:  100,
			},
			mockSetup: func(m *MockFinnhubClient) {
				m.AddQuote("TSLA", &QuoteResponse{
					CurrentPrice: 100,
				})
			},
			wantValue:  0,
			wantProfit: 0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設置 mock
			mock := NewMockFinnhubClient()
			tt.mockSetup(mock)
			SetDefaultClient(mock)
			defer ResetDefaultClient()

			gotValue, gotProfit, err := Calculate(ctx, tt.input)

			// 驗證錯誤
			if tt.wantErr {
				if err == nil {
					t.Errorf("Calculate() error = nil, wantErr = true")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Calculate() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Calculate() unexpected error = %v", err)
				return
			}

			// 驗證結果（使用小誤差容忍）
			epsilon := 0.01
			if diff := gotValue - tt.wantValue; diff < -epsilon || diff > epsilon {
				t.Errorf("Calculate() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if diff := gotProfit - tt.wantProfit; diff < -epsilon || diff > epsilon {
				t.Errorf("Calculate() gotProfit = %v, want %v", gotProfit, tt.wantProfit)
			}
		})
	}
}
