package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"discordBot/service/discord"
	"discordBot/service/stock"
)

// Quote : 取得股價
// example : $+TSLA
func Quote(s *discordgo.Session, m *discordgo.MessageCreate) {
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

// SetStock : 新增股票到 DB
// example : $set_stock TSLA units price
func SetStock(s *discordgo.Session, m *discordgo.MessageCreate) {
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

// GetStock : 取得 DB 中股票
// example : $get_stock TSLA
func GetStock(s *discordgo.Session, m *discordgo.MessageCreate) {
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
