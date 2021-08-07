package handler

import (
	"discordBot/service/stock"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Quote(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$+") {
		res, err := stock.Quote(m.Content)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error:%v", err))
		}
		s.ChannelMessageSend(m.ChannelID, res)
	}

}

func Cron(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "$開始背景") {
		s.ChannelMessageSend(m.ChannelID, "開始背景")

		for {
			res, err := stock.GetChange("ARKK")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error:%v", err))
			}
			if res > 1 {
				s.ChannelMessageSend(m.ChannelID, "ARKK 漲幅超過1%, 賺錢囉")
				time.Sleep(time.Second * 10)
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}
}
