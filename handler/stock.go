package handler

import (
	"context"
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
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("錯誤: %v", err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, res)
	}

}

// SetStock : 新增股票到 DB
func SetStock(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// example : $set_stock TSLA units price
	if strings.HasPrefix(m.Content, "$set_stock") {
		err := stock.SetStock(m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("錯誤: %v", err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, "新增成功")
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
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("錯誤: %v", err))
			return
		}

		marshalRes, err := json.Marshal(res)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("錯誤: %v", err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, string(marshalRes))
	}
}
