package stock

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"discordBot/model/redis"
	"discordBot/service/discord"
)

// CheckChange : 檢查漲跌幅
func CheckChange(s *discordgo.Session) {
	ctx := context.Background()

	// Redis 取出資料
	watchList, err := redis.LRange(
		ctx,
		"watch_list",
		0,
		-1,
	)
	if err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: "872317320729616395",
				Content:   fmt.Sprintf("取得列表時錯誤: %v", err),
			},
		)
		return
	}

	for _, v := range watchList {
		// 先看是否已通知過
		redisRes, err := redis.Get(ctx, "watch_list:"+v)
		if err != nil {
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: "872317320729616395",
					Content:   fmt.Sprintf("取得紀錄時錯誤: %v", err),
				},
			)
			return
		}

		if redisRes != "" {
			continue
		}

		change, err := GetChange(ctx, v)
		if err != nil {
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: "872317320729616395",
					Content:   fmt.Sprintf("取得漲跌幅時錯誤: %v", err),
				},
			)
			continue
		}

		if change > 3 || change < -3 {
			_, err = s.ChannelMessageSendComplex("872317320729616395", &discordgo.MessageSend{
				Content: fmt.Sprintf("<@512265930735222795> 警告: %s 今日漲跌幅為 %.2f %s", v, change, "%"),
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
				},
			})
			if err != nil {
				discord.SendMessage(
					s,
					&discord.SendMessageInput{
						ChannelID: "872317320729616395",
						Content:   fmt.Sprintf("發送訊息時錯誤: %v", err),
					},
				)
				return
			}

			// 寫入紀錄已通知
			err = redis.Set(ctx, "watch_list:"+v, "true", time.Hour*8)
			if err != nil {
				discord.SendMessage(
					s,
					&discord.SendMessageInput{
						ChannelID: "872317320729616395",
						Content:   fmt.Sprintf("寫入紀錄時錯誤: %v", err),
					},
				)
				return
			}
		}
	}
}
