package handler

import (
	"context"
	"discordBot/model/redis"
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

		key := strSlice[1]
		value := strSlice[2]

		err := redis.Set(
			context.Background(),
			key,
			value,
		)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error:%v", err))
		}

		res := fmt.Sprintf("設定 key: %s, value: %s", key, value)
		s.ChannelMessageSend(m.ChannelID, res)
	}
}

func GetRedis(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$getRedis") {
		strSlice := strings.Split(m.Content, " ")

		key := strSlice[1]

		value, err := redis.Get(
			context.Background(),
			key,
		)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error:%v", err))
		}

		res := fmt.Sprintf("取得 key: %s, value: %s", key, value)

		s.ChannelMessageSend(m.ChannelID, res)
	}
}
