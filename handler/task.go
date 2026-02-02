package handler

import (
	"fmt"
	"time"

	"discordBot/pkg/config"
	"discordBot/pkg/logger"
	"discordBot/service/crypto"
	"discordBot/service/discord"
	"discordBot/service/stock"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

var taskCron *cron.Cron

// Task : 定時任務
func Task(s *discordgo.Session) {
	taskConfig := config.GetTaskConfig()

	// 設定時區為台北
	taipeiLoc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		logger.Error("無法載入台北時區", "error", err)
		taipeiLoc = time.UTC
	}

	// 新建一個定時任務物件
	taskCron = cron.New(cron.WithLocation(taipeiLoc))

	// 週一到週五 23:00 - 23:59 每 10 分鐘啟動
	taskCron.AddFunc("*/10 23 * * 1-5", func() {
		logger.Info("執行股票漲跌幅檢查任務")
		stock.CheckChange(s)
	})

	// 週二到週六 00:00 - 04 : 59 每 10 分鐘啟動
	taskCron.AddFunc("*/10 0-4 * * 2-6", func() {
		logger.Info("執行股票漲跌幅檢查任務")
		stock.CheckChange(s)
	})

	// 週二到週六 06:00 啟動
	taskCron.AddFunc("0 6 * * 2-6", func() {
		logger.Info("執行收益計算任務")
		stock.CalculateProfit(s)
	})

	taskCron.AddFunc("@every 30s", func() {
		price, err := crypto.GetPrice()
		if err != nil {
			logger.Error("取得ETH價格失敗", "error", err)
			discord.SendMessage(
				s,
				&discord.SendMessageInput{
					ChannelID: taskConfig.CryptoPriceChannelID,
					Content:   fmt.Sprintf("ETH 取得價格錯誤: %v", err),
				},
			)
		}
		s.UpdateListeningStatus(fmt.Sprintf("ETH價格 %.2F", price))
	})

	taskCron.Start()
	logger.Info("定時任務已啟動")
}

// StopTasks 停止所有定時任務
func StopTasks() {
	if taskCron != nil {
		taskCron.Stop()
		logger.Info("定時任務已停止")
	}
}
