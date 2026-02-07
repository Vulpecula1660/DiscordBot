package handler

import (
	"context"
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

var (
	taskCron          *cron.Cron
	cronErrorReporter = discord.NewErrorReporter(time.Minute)
)

// Task : 定時任務
func Task(s *discordgo.Session) {
	taskConfig := config.GetTaskConfig()
	cronErrorReporter.SetCooldown(taskDurationFromSeconds(taskConfig.ErrorNotifyCooldownSeconds, time.Minute))
	externalTimeout := taskDurationFromSeconds(taskConfig.ExternalCallTimeoutSeconds, 15*time.Second)

	// 設定時區為台北
	taipeiLoc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		logger.Error("無法載入台北時區", "error", err)
		taipeiLoc = time.UTC
	}

	// 新建一個定時任務物件
	taskCron = cron.New(cron.WithLocation(taipeiLoc))
	registeredCount := 0

	// 週一到週五 23:00 - 23:59 每 10 分鐘啟動
	if registerTask("*/10 23 * * 1-5", "check_change_weekday_late", func() {
		logger.Info("執行股票漲跌幅檢查任務")
		stock.CheckChange(s)
	}) {
		registeredCount++
	}

	// 週二到週六 00:00 - 04 : 59 每 10 分鐘啟動
	if registerTask("*/10 0-4 * * 2-6", "check_change_weekday_early", func() {
		logger.Info("執行股票漲跌幅檢查任務")
		stock.CheckChange(s)
	}) {
		registeredCount++
	}

	// 週二到週六 06:00 啟動
	if registerTask("0 6 * * 2-6", "calculate_profit_daily", func() {
		logger.Info("執行收益計算任務")
		stock.CalculateProfit(s)
	}) {
		registeredCount++
	}

	if registerTask("@every 30s", "crypto_price_update", func() {
		priceCtx, priceCancel := context.WithTimeout(context.Background(), externalTimeout)
		price, err := crypto.GetPriceWithContext(priceCtx)
		priceCancel()
		if err != nil {
			logger.Error("取得ETH價格失敗", "error", err)
			cronErrorReporter.Notify(
				s,
				"task:crypto:get_price",
				&discord.SendMessageInput{
					ChannelID: taskConfig.CryptoPriceChannelID,
					Content:   fmt.Sprintf("ETH 取得價格錯誤: %v", err),
				},
			)
			return
		}

		if err := s.UpdateListeningStatus(fmt.Sprintf("ETH價格 %.2F", price)); err != nil {
			logger.Error("更新 Discord 狀態失敗", "error", err)
		}
	}) {
		registeredCount++
	}

	if registeredCount == 0 {
		logger.Error("沒有成功註冊任何定時任務")
		return
	}

	taskCron.Start()
	logger.Info("定時任務已啟動", "registeredCount", registeredCount)
}

func registerTask(spec string, name string, fn func()) bool {
	if _, err := taskCron.AddFunc(spec, fn); err != nil {
		logger.Error("註冊定時任務失敗", "task", name, "spec", spec, "error", err)
		return false
	}

	logger.Info("註冊定時任務成功", "task", name, "spec", spec)
	return true
}

func taskDurationFromSeconds(seconds int, defaultDuration time.Duration) time.Duration {
	if seconds <= 0 {
		return defaultDuration
	}

	return time.Duration(seconds) * time.Second
}

// StopTasks 停止所有定時任務
func StopTasks() {
	if taskCron != nil {
		taskCron.Stop()
		logger.Info("定時任務已停止")
	}
}
