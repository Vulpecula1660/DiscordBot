package handler

import (
	"context"
	"discordBot/service/discord"
	"discordBot/service/stock"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Quote : 取得股價
func Quote(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// example : $+TSLA
	if strings.HasPrefix(m.Content, "$+") {
		res, err := stock.Quote(context.Background(), m.Content)
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

		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   res,
			},
		)
	}

}

// SetStock : 新增股票到 DB
func SetStock(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// example : $set_stock TSLA units price
	if strings.HasPrefix(m.Content, "$set_stock") {
		err := stock.SetStock(context.Background(), m)
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

		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   "新增成功",
			},
		)
	}
}

// 取得 DB 中股票
func GetStock(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// example : $get_stock TSLA
	if strings.HasPrefix(m.Content, "$get_stock") {
		res, err := stock.GetStock(context.Background(), m)
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

		marshalRes, err := json.Marshal(res)
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

		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: m.ChannelID,
				Content:   string(marshalRes),
			},
		)
	}
}
