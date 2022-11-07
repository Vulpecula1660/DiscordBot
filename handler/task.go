package handler

import (
	"fmt"
	"time"

	"discordBot/service/crypto"
	"discordBot/service/discord"
	"discordBot/service/stock"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// Task : 定時任務
func Task(s *discordgo.Session) {
	// 設定時區為台北
	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")

	// 新建一個定時任務物件
	c := cron.New(cron.WithLocation(taipeiLoc))

	// 週一到週五 23:00 - 23:59 每 10 分鐘啟動
	c.AddFunc("*/10 23 * * 1-5", func() {
		// 取得漲跌幅
		stock.CheckChange(s)
	})

	// 週二到週六 00:00 - 04 : 59 每 10 分鐘啟動
	c.AddFunc("*/10 0-4 * * 2-6", func() {
		// 取得漲跌幅
		stock.CheckChange(s)
	})

	// 週二到週六 06:00 啟動
	c.AddFunc("0 6 * * 2-6", func() {
		// 計算收益
		stock.CalculateProfit(s)
	})

	c.AddFunc("@every 30s", func() {
		price, err := crypto.GetPrice()
		if err != nil {
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: "1032641300077490266",
					Content:   fmt.Sprintf("ETH 取得價格錯誤: %v", err),
				},
			)
		}
		s.UpdateListeningStatus(fmt.Sprintf("ETH價格 %.2F", price))
	})

	c.Start()
}
