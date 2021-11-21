package stock

import (
	"context"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"

	"discordBot/model/redis"
)

func Test_CheckChange(t *testing.T) {
	// 使用 Monkey Patch
	patchLRange := monkey.Patch(redis.LRange, func(ctx context.Context, key string, start int64, stop int64) ([]string, error) {
		return []string{"TSLA", "AAPL"}, nil
	})

	patchGet := monkey.Patch(redis.Get, func(ctx context.Context, key string) (string, error) {
		return "", nil
	})

	patchGetChange := monkey.Patch(GetChange, func(ctx context.Context, stock string) (float32, error) {
		return 5, nil
	})

	patchSet := monkey.Patch(redis.Set, func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
		return nil
	})

	defer func() {
		patchLRange.Restore()
		patchGet.Restore()
		patchGetChange.Restore()
		patchSet.Restore()
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
			CheckChange(tt.args.s)
		})
	}
}
