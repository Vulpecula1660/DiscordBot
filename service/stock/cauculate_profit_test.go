package stock

import (
	"context"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"

	"discordBot/model/dao/stock"
	"discordBot/model/dto"
	"discordBot/model/redis"
)

func TestCalculateProfit(t *testing.T) {
	// 使用 Monkey Patch
	patchGet := monkey.Patch(stock.Get, func(ctx context.Context, input *stock.GetInput) (ret []*dto.Stock, err error) {
		return []*dto.Stock{
			{
				ID:     0,
				UserID: "",
				Symbol: "TSLA",
				Units:  1,
				Price:  1,
			},
			{
				ID:     0,
				UserID: "",
				Symbol: "AAPL",
				Units:  1,
				Price:  1,
			}}, nil
	})

	patchCalculate := monkey.Patch(Calculate, func(ctx context.Context, input *CalculateInput) (value float64, profit float64, err error) {
		return 1000, 999, nil
	})

	patchRedisGet := monkey.Patch(redis.Get, func(ctx context.Context, key string) (string, error) {
		return "1", nil
	})

	patchRedisSet := monkey.Patch(redis.Set, func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
		return nil
	})

	defer func() {
		patchGet.Restore()
		patchCalculate.Restore()
		patchRedisGet.Restore()
		patchRedisSet.Restore()
	}()

	token := os.Getenv("DCToken")

	// creates a new Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		t.Error("error creating Discord session,", err)
		return
	}

	type args struct {
		s *discordgo.Session
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success Test",
			args: args{
				s: dg,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CalculateProfit(tt.args.s)
		})
	}
}
