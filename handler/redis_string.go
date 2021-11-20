package handler

import (
	"context"
	"discordBot/model/redis"
	"discordBot/service/discord"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SetRedis(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$setRedis") {
		strSlice := strings.Split(m.Content, " ")

		if len(strSlice) != 3 {
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: m.ChannelID,
					Content:   "參數錯誤",
				},
			)

			return
		}

		key := strSlice[1]
		value := strSlice[2]

		err := redis.Set(
			context.Background(),
			key,
			value,
			0, // 無限時
		)
		if err != nil {
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: m.ChannelID,
					Content:   fmt.Sprintf("錯誤: %v", err),
				},
			)
			return
		}

		res := fmt.Sprintf("設定 key: %s, value: %s", key, value)

		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   res,
			},
		)
	}
}

func GetRedis(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$getRedis") {
		strSlice := strings.Split(m.Content, " ")

		if len(strSlice) != 2 {
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: m.ChannelID,
					Content:   "參數錯誤",
				},
			)
			return
		}

		key := strSlice[1]

		value, err := redis.Get(
			context.Background(),
			key,
		)
		if err != nil {
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: m.ChannelID,
					Content:   fmt.Sprintf("錯誤: %v", err),
				},
			)
			return
		}

		res := fmt.Sprintf("取得 key: %s, value: %s", key, value)

		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   res,
			},
		)
	}
}
