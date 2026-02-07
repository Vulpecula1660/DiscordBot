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

func SetList(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 3 {
		if err := discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $setList <key> <value>",
			},
		); err != nil {
			logger.Error("發送訊息失敗", "error", err)
		}
		return
	}

	key := strSlice[1]
	value := strSlice[2]

	err := redis.RPush(
		context.Background(),
		key,
		value,
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

func GetList(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 2 {
		if err := discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $getList <key>",
			},
		); err != nil {
			logger.Error("發送訊息失敗", "error", err)
		}
		return
	}

	key := strSlice[1]

	value, err := redis.LRange(
		context.Background(),
		key,
		0,
		-1,
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

func DelListValue(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 3 {
		if err := discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $delListValue <key> <value>",
			},
		); err != nil {
			logger.Error("發送訊息失敗", "error", err)
		}
		return
	}

	key := strSlice[1]
	value := strSlice[2]

	err := redis.LRem(
		context.Background(),
		key,
		0,
		value,
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

	res := fmt.Sprintf("從 key: %s 中刪除 value: %s", key, value)

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
