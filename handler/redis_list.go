package handler

import (
	"context"
	"discordBot/model/redis"
	"discordBot/service/discord"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SetList(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$setList") {
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
}

func GetList(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$getList") {
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
}

func DelListValue(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$delListValue") {
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
}
