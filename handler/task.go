package handler

import (
	"discordBot/service/stock"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

func Task(s *discordgo.Session) {
	// 設定時區為台北
	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")

	// 新建一個定時任務物件
	c := cron.New(cron.WithLocation(taipeiLoc))

	c.AddFunc("*/10 22-23 * * 1-5", func() {
		stock.CheckChange(s)
	})

	c.AddFunc("*/10 0-4 * * 2-6", func() {
		stock.CheckChange(s)
	})

	c.AddFunc("0 5 * * 2-6", func() {
		stock.CalculateProfit(s)
	})

	c.Start()
}
