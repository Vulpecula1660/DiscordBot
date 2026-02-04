package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordBot/model/redis"
	"discordBot/pkg/validator"
	"discordBot/service/discord"
)

// SetRedis : 設定 Redis 字串值
// example : $setRedis key value
func SetRedis(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	key := validator.SanitizeInput(strSlice[1])
	value := validator.SanitizeInput(strSlice[2])

	// 驗證 key
	if err := validator.ValidateRedisKey(key); err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   fmt.Sprintf("Key 驗證失敗: %v", err),
			},
		)
		return
	}

	// 驗證 value
	if err := validator.ValidateRedisValue(value); err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   fmt.Sprintf("Value 驗證失敗: %v", err),
			},
		)
		return
	}

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

// GetRedis : 取得 Redis 字串值
// example : $getRedis key
func GetRedis(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	key := validator.SanitizeInput(strSlice[1])

	// 驗證 key
	if err := validator.ValidateRedisKey(key); err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   fmt.Sprintf("Key 驗證失敗: %v", err),
			},
		)
		return
	}

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
