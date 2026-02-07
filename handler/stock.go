package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"discordBot/pkg/logger"
	"discordBot/service/discord"
	"discordBot/service/stock"
)

// Quote : 取得股價
func Quote(s *discordgo.Session, m *discordgo.MessageCreate) {
	// example : $+TSLA
	res, err := stock.Quote(context.Background(), m.Content)
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

// SetStock : 新增股票到 DB
func SetStock(s *discordgo.Session, m *discordgo.MessageCreate) {
	// example : $set_stock TSLA units price
	err := stock.SetStock(context.Background(), m)
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

	if err := discord.SendMessage(
		s,
		&discord.SendMessageInput{
			ChannelID: m.ChannelID,
			Content:   "新增成功",
		},
	); err != nil {
		logger.Error("發送訊息失敗", "error", err)
	}
}

// GetStock : 取得 DB 中股票
func GetStock(s *discordgo.Session, m *discordgo.MessageCreate) {
	// example : $get_stock TSLA
	res, err := stock.GetStock(context.Background(), m)
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

	marshalRes, err := json.Marshal(res)
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

	if err := discord.SendMessage(
		s,
		&discord.SendMessageInput{
			ChannelID: m.ChannelID,
			Content:   string(marshalRes),
		},
	); err != nil {
		logger.Error("發送訊息失敗", "error", err)
	}
}
