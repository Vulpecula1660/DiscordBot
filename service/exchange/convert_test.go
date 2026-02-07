package exchange

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"
)

func Test_ConvertExchange(t *testing.T) {
	tests := []struct {
		name         string
		oldMoney     []float64
		mockRate     float64
		mockErr      error
		wantNewMoney []float64
		wantErr      bool
		errContains  string
	}{
		{
			name:         "successful conversion - USD to TWD",
			oldMoney:     []float64{1, 10, 50},
			mockRate:     27.7995,
			wantNewMoney: []float64{27.7995, 277.995, 1389.975},
			wantErr:      false,
		},
		{
			name:         "empty input",
			oldMoney:     []float64{},
			mockRate:     30.0,
			wantNewMoney: []float64{},
			wantErr:      false,
		},
		{
			name:         "single value",
			oldMoney:     []float64{100},
			mockRate:     30.5,
			wantNewMoney: []float64{3050},
			wantErr:      false,
		},
		{
			name:         "zero values",
			oldMoney:     []float64{0, 0, 0},
			mockRate:     27.5,
			wantNewMoney: []float64{0, 0, 0},
			wantErr:      false,
		},
		{
			name:        "API error",
			oldMoney:    []float64{1, 2, 3},
			mockErr:     errors.New("connection timeout"),
			wantErr:     true,
			errContains: "connection timeout",
		},
		{
			name:         "decimal precision",
			oldMoney:     []float64{0.1, 0.01, 0.001},
			mockRate:     30.123,
			wantNewMoney: []float64{3.0123, 0.30123, 0.030123},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設置 mock
			mock := &MockRateProvider{
				Rate: tt.mockRate,
				Err:  tt.mockErr,
			}
			SetRateProvider(mock)
			defer ResetRateProvider()

			gotNewMoney, err := ConvertExchange(tt.oldMoney)

			// 驗證錯誤
			if tt.wantErr {
				if err == nil {
					t.Errorf("ConvertExchange() error = nil, wantErr = true")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ConvertExchange() error = %v, want error containing %v", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ConvertExchange() unexpected error = %v", err)
				return
			}

			// 驗證結果（使用小誤差容忍）
			epsilon := 0.0001
			if !floatSlicesEqual(gotNewMoney, tt.wantNewMoney, epsilon) {
				t.Errorf("ConvertExchange() = %v, want %v", gotNewMoney, tt.wantNewMoney)
			}
		})
	}
}

func Test_ConvertExchangeWithContext(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		oldMoney     []float64
		mockRate     float64
		mockErr      error
		wantNewMoney []float64
		wantErr      bool
	}{
		{
			name:         "successful with context",
			oldMoney:     []float64{100},
			mockRate:     30.0,
			wantNewMoney: []float64{3000},
			wantErr:      false,
		},
		{
			name:     "context cancellation",
			oldMoney: []float64{1},
			mockErr:  errors.New("context cancelled"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設置 mock
			mock := &MockRateProvider{
				Rate: tt.mockRate,
				Err:  tt.mockErr,
			}
			SetRateProvider(mock)
			defer ResetRateProvider()

			gotNewMoney, err := ConvertExchangeWithContext(ctx, tt.oldMoney)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ConvertExchangeWithContext() error = nil, wantErr = true")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("ConvertExchangeWithContext() unexpected error = %v", err)
				return
			}

			if !reflect.DeepEqual(gotNewMoney, tt.wantNewMoney) {
				t.Errorf("ConvertExchangeWithContext() = %v, want %v", gotNewMoney, tt.wantNewMoney)
			}
		})
	}
}

// 輔助函數
func contains(s, substr string) bool {
	return len(substr) <= len(s) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func floatSlicesEqual(a, b []float64, epsilon float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > epsilon {
			return false
		}
	}
	return true
}
