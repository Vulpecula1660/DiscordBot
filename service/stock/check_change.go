package stock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"discordBot/model/redis"
	"discordBot/pkg/config"
	"discordBot/pkg/logger"
	"discordBot/service/discord"

	"github.com/bwmarrin/discordgo"
)

// CheckChange : 檢查漲跌幅
func CheckChange(s *discordgo.Session) {
	ctx := context.Background()
	taskConfig := config.GetTaskConfig()

	logger.Info("開始檢查股票漲跌幅")

	// Redis 取出資料
	watchList, err := redis.LRange(
		ctx,
		"watch_list",
		0,
		-1,
	)
	if err != nil {
		logger.Error("取得觀察列表失敗", "error", err)
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: taskConfig.WatchListChannelID,
				Content:   fmt.Sprintf("取得列表時錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("取得觀察列表", "count", len(watchList))

	var wg sync.WaitGroup
	wg.Add(len(watchList))

	for _, v := range watchList {
		go func(symbol string) {
			defer wg.Done()
			// 先看是否已通知過
			redisRes, err := redis.Get(ctx, "watch_list:"+symbol)
			if err != nil {
				logger.Error("取得通知紀錄失敗", "symbol", symbol, "error", err)
				discord.SendMessage(
					s,
					&discord.SendMessageInput{
						ChannelID: taskConfig.WatchListChannelID,
						Content:   fmt.Sprintf("取得紀錄時錯誤: %v", err),
					},
				)
				return
			}

			if redisRes != "" {
				return
			}

			change, err := GetChange(ctx, symbol)
			if err != nil {
				logger.Error("取得漲跌幅失敗", "symbol", symbol, "error", err)
				discord.SendMessage(
					s,
					&discord.SendMessageInput{
						ChannelID: taskConfig.WatchListChannelID,
						Content:   fmt.Sprintf("取得漲跌幅時錯誤: %v", err),
					},
				)
				return
			}

			if change > 3 || change < -3 {
				logger.Warn("股票漲跌幅超過閾值", "symbol", symbol, "change", change)
				_, err = s.ChannelMessageSendComplex(taskConfig.WatchListChannelID, &discordgo.MessageSend{
					Content: fmt.Sprintf("<@%s> 警告: %s 今日漲跌幅為 %.2f %%", taskConfig.DefaultUserID, symbol, change),
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
					},
				})
				if err != nil {
					logger.Error("發送警告訊息失敗", "symbol", symbol, "error", err)
					discord.SendMessage(
						s,
						&discord.SendMessageInput{
							ChannelID: taskConfig.WatchListChannelID,
							Content:   fmt.Sprintf("發送訊息時錯誤: %v", err),
						},
					)
					return
				}

				// 寫入紀錄已通知
				err = redis.Set(ctx, "watch_list:"+symbol, "true", time.Hour*8)
				if err != nil {
					logger.Error("寫入通知紀錄失敗", "symbol", symbol, "error", err)
					discord.SendMessage(
						s,
						&discord.SendMessageInput{
							ChannelID: taskConfig.WatchListChannelID,
							Content:   fmt.Sprintf("寫入紀錄時錯誤: %v", err),
						},
					)
					return
				}
			}
		}(v)
	}

	wg.Wait()
	logger.Info("完成股票漲跌幅檢查")
}
