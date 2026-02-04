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

// SetList : 新增值到 Redis List
// example : $setList key value
func SetList(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	err := redis.RPush(
		context.Background(),
		key,
		value,
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

// GetList : 取得 Redis List 所有值
// example : $getList key
func GetList(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	value, err := redis.LRange(
		context.Background(),
		key,
		0,
		-1,
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

// DelListValue : 從 Redis List 刪除指定值
// example : $delListValue key value
func DelListValue(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	err := redis.LRem(
		context.Background(),
		key,
		0,
		value,
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

	res := fmt.Sprintf("從 key: %s 中刪除 value: %s", key, value)

	discord.SendMessage(
		s,
		&discord.SendMessageInput{
			ChannelID: m.ChannelID,
			Content:   res,
		},
	)
}
