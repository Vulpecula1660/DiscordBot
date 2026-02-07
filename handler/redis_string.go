package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordBot/model/redis"
	"discordBot/pkg/logger"
	"discordBot/service/discord"
)

func SetRedis(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 3 {
		if err := discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $setRedis <key> <value>",
			},
		); err != nil {
			logger.Error("發送訊息失敗", "error", err)
		}
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
		if err := discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   fmt.Sprintf("錯誤: %v", err),
			},
		); err != nil {
			logger.Error("發送訊息失敗", "error", err)
		}
		return
	}

	res := fmt.Sprintf("設定 key: %s, value: %s", key, value)

	if err := discord.SendMessage(
		s,
		&discord.SendMessageInput{
			ChannelID: m.ChannelID,
			Content:   res,
		},
	); err != nil {
		logger.Error("發送訊息失敗", "error", err)
	}
}

func GetRedis(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 2 {
		if err := discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $getRedis <key>",
			},
		); err != nil {
			logger.Error("發送訊息失敗", "error", err)
		}
		return
	}

	key := strSlice[1]

	value, err := redis.Get(
		context.Background(),
		key,
	)
	if err != nil {
		if err := discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   fmt.Sprintf("錯誤: %v", err),
			},
		); err != nil {
			logger.Error("發送訊息失敗", "error", err)
		}
		return
	}

	res := fmt.Sprintf("取得 key: %s, value: %s", key, value)

	if err := discord.SendMessage(
		s,
		&discord.SendMessageInput{
			ChannelID: m.ChannelID,
			Content:   res,
		},
	); err != nil {
		logger.Error("發送訊息失敗", "error", err)
	}
}
