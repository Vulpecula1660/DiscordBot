package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordBot/model/redis"
	"discordBot/service/discord"
)

func SetList(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 3 {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $setList <key> <value>",
			},
		)
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

func GetList(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 2 {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $getList <key>",
			},
		)
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

func DelListValue(s *discordgo.Session, m *discordgo.MessageCreate) {
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 3 {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "參數錯誤，格式: $delListValue <key> <value>",
			},
		)
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
