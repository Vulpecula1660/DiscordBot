package stock

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func Test_Quote(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		err     error
	}{
		{
			name: "TSLA",
			args: args{
				ctx:     ctx,
				message: "$+TSLA",
			},
			want:    "Finnhub 查詢標的為:TSLA 目前價格為:1137.06 今天漲跌幅:3.7104%",
			wantErr: false,
			err:     nil,
		},
		{
			name: "Wrong Message",
			args: args{
				ctx:     ctx,
				message: "aaaaa",
			},
			want:    "",
			wantErr: true,
			err:     fmt.Errorf("參數錯誤"),
		},
		{
			name: "Wrong Symbol",
			args: args{
				ctx:     ctx,
				message: "$+aaaaa",
			},
			want:    "",
			wantErr: true,
			err:     fmt.Errorf("搜尋失敗"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Quote(tt.args.ctx, tt.args.message)
			if err != nil && tt.wantErr {
				if err.Error() != tt.err.Error() {
					t.Errorf("Quote() error = %v, wantErr %v", err.Error(), tt.err.Error())
					return
				}
			}
			if got != tt.want {
				t.Errorf("Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetChange(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		stock string
	}
	tests := []struct {
		name    string
		args    args
		want    float32
		wantErr bool
		err     error
	}{
		{
			name: "TSLA",
			args: args{
				ctx:   ctx,
				stock: "Tsla",
			},
			want:    3.7104,
			wantErr: false,
			err:     nil,
		},
		{
			name: "Wrong Symbol",
			args: args{
				ctx:   ctx,
				stock: "01346",
			},
			want:    0,
			wantErr: true,
			err:     fmt.Errorf("搜尋失敗"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChange(tt.args.ctx, tt.args.stock)
			if err != nil && tt.wantErr {
				if err.Error() != tt.err.Error() {
					t.Errorf("Quote() error = %v, wantErr %v", err.Error(), tt.err.Error())
					return
				}
			}
			if got != tt.want {
				t.Errorf("GetChange() = %v, want %v, err = %v", got, tt.want, err)
			}
		})
	}
}

func Test_Calculate(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		input *CalculateInput
	}
	tests := []struct {
		name       string
		args       args
		wantValue  float64
		wantProfit float64
		wantErr    bool
		err        error
	}{
		{
			name: "TSLA",
			args: args{
				ctx: ctx,
				input: &CalculateInput{
					Symbol: "TSLA",
					Units:  1,
					Price:  1,
				},
			},
			wantValue:  1137.06005859375,
			wantProfit: 1136.06005859375,
			wantErr:    false,
			err:        nil,
		},
		{
			name: "Wrong Symbol",
			args: args{
				ctx: ctx,
				input: &CalculateInput{
					Symbol: "012346",
					Units:  1,
					Price:  1,
				},
			},
			wantValue:  0,
			wantProfit: 0,
			wantErr:    true,
			err:        fmt.Errorf("搜尋失敗"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotProfit, err := Calculate(tt.args.ctx, tt.args.input)
			if err != nil && tt.wantErr {
				if err.Error() != tt.err.Error() {
					t.Errorf("Quote() error = %v, wantErr %v", err.Error(), tt.err.Error())
					return
				}
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
