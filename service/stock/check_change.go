package stock

import (
	"context"
	"discordBot/model/redis"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CheckChange(s *discordgo.Session) {
	ctx := context.Background()

	watchList, err := redis.LRange(
		ctx,
		"watch_list",
		0,
		-1,
	)
	if err != nil {
		s.ChannelMessageSend("872317320729616395", fmt.Sprintf("取得列表時錯誤: %v", err))
		return
	}

	for _, v := range watchList {
		change, err := GetChange(v)
		if err != nil {
			s.ChannelMessageSend("872317320729616395", fmt.Sprintf("取得漲跌幅時錯誤: %v", err))
			continue
		}

		if change > 3 || change < -3 {
			// 先看是否已通知過
			res, err := redis.Get(ctx, "watch_list:"+v)
			if err != nil {
				s.ChannelMessageSend("872317320729616395", fmt.Sprintf("取得紀錄時錯誤: %v", err))
				return
			}

			if res != "" {
				continue
			}

			_, err = s.ChannelMessageSendComplex("872317320729616395", &discordgo.MessageSend{
				Content: fmt.Sprintf("<@512265930735222795> 警告: %s 今日漲跌幅為 %v %s", v, change, "%"),
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
				},
			})
			if err != nil {
				s.ChannelMessageSend("872317320729616395", fmt.Sprintf("發送訊息時錯誤: %v", err))
				return
			}

			// 寫入紀錄已通知
			err = redis.Set(ctx, "watch_list:"+v, "true", time.Hour*8)
			if err != nil {
				s.ChannelMessageSend("872317320729616395", fmt.Sprintf("寫入紀錄時錯誤: %v", err))
				return
			}
		}
	}
}
