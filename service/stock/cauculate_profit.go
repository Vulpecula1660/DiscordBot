package stock

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	stockdao "discordBot/model/dao/stock"
	"discordBot/model/dto"
	"discordBot/pkg/config"
	"discordBot/pkg/logger"
	"discordBot/service/discord"
	"discordBot/service/exchange"

	"github.com/bwmarrin/discordgo"
)

// StockRepository 股票數據倉庫接口類型
type StockRepository interface {
	Get(ctx context.Context, input *stockdao.GetInput) ([]*dto.Stock, error)
}

// CalculateProfit : 計算損益
func CalculateProfit(s *discordgo.Session) {
	CalculateProfitWithDeps(s, stockDaoDeps{}, redisDeps{})
}

// stockDaoDeps 封裝 Stock DAO 依賴
type stockDaoDeps struct{}

func (d stockDaoDeps) Get(ctx context.Context, input *stockdao.GetInput) ([]*dto.Stock, error) {
	return stockdao.Get(ctx, input)
}

// CalculateProfitWithDeps 使用指定依賴計算損益（用於測試）
func CalculateProfitWithDeps(s *discordgo.Session, repo StockRepository, redisClient RedisClient) {
	taskConfig := config.GetTaskConfig()
	taskErrorReporter.SetCooldown(durationFromSeconds(taskConfig.ErrorNotifyCooldownSeconds, time.Minute))

	runTimeout := durationFromSeconds(taskConfig.CalculateProfitTimeoutSeconds, 3*time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()

	externalTimeout := durationFromSeconds(taskConfig.ExternalCallTimeoutSeconds, 15*time.Second)
	maxConcurrency := normalizeConcurrency(taskConfig.CalculateProfitMaxConcurrency, 5)

	logger.Info("開始計算收益", "maxConcurrency", maxConcurrency, "timeout", runTimeout.String())

	// 先取資料
	getStockCtx, getStockCancel := context.WithTimeout(ctx, externalTimeout)
	dbRes, err := repo.Get(
		getStockCtx,
		&stockdao.GetInput{
			UserID: taskConfig.DefaultUserID,
		},
	)
	getStockCancel()
	if err != nil {
		logger.Error("取得股票資料失敗", "error", err)
		taskErrorReporter.Notify(
			s,
			"calculate_profit:get_stock",
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("取資料時錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("取得股票資料", "count", len(dbRes))
	if len(dbRes) == 0 {
		logger.Info("持倉為空，略過收益計算")
		return
	}

	var totalProfit, totalCost, totalValue float64

	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, maxConcurrency)

loop:
	for _, v := range dbRes {
		if err := ctx.Err(); err != nil {
			logger.Warn("收益計算任務超時，停止派發剩餘標的", "error", err)
			break
		}

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			logger.Warn("收益計算任務超時，停止派發剩餘標的", "error", ctx.Err())
			break loop
		}

		wg.Add(1)
		go func(holding *dto.Stock) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			calcCtx, calcCancel := context.WithTimeout(ctx, externalTimeout)
			value, profit, err := Calculate(
				calcCtx,
				&CalculateInput{
					Symbol: holding.Symbol,
					Units:  holding.Units,
					Price:  holding.Price,
				})
			calcCancel()
			if err != nil {
				logger.Error("計算損益失敗", "symbol", holding.Symbol, "error", err)
				taskErrorReporter.Notify(
					s,
					"calculate_profit:calculate",
					&discord.SendMessageInput{
						ChannelID: taskConfig.ProfitReportChannelID,
						Content:   fmt.Sprintf("計算損益時錯誤: %v", err),
					},
				)
				return
			}

			mu.Lock()
			cost := holding.Units * holding.Price

			totalProfit += profit
			totalCost += cost
			totalValue += value
			mu.Unlock()
		}(v)
	}

	wg.Wait()
	if err := ctx.Err(); err != nil && err != context.Canceled {
		logger.Warn("收益計算任務逾時", "error", err)
		taskErrorReporter.Notify(
			s,
			"calculate_profit:timeout",
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("收益計算任務逾時: %v", err),
			},
		)
		return
	}

	// 取昨日市場總值
	redisKey := taskConfig.ProfitReportChannelID + "_" + "totalValue"
	redisGetCtx, redisGetCancel := context.WithTimeout(ctx, externalTimeout)
	yesterdayTotalValue, err := redisClient.Get(redisGetCtx, redisKey)
	redisGetCancel()
	if err != nil {
		logger.Error("取得昨日市場總值失敗", "error", err)
		taskErrorReporter.Notify(
			s,
			"calculate_profit:redis_get",
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("取昨日市場總值錯誤: %v", err),
			},
		)
		return
	}

	var yesterdayTotalValueFloat float64
	if yesterdayTotalValue == "" {
		yesterdayTotalValueFloat = totalValue
		logger.Info("昨日市場總值不存在，使用今日總值作為基準", "redisKey", redisKey)
	} else {
		yesterdayTotalValueFloat, err = strconv.ParseFloat(yesterdayTotalValue, 64)
		if err != nil {
			logger.Error("轉換昨日市場總值失敗", "value", yesterdayTotalValue, "error", err)
			taskErrorReporter.Notify(
				s,
				"calculate_profit:redis_parse",
				&discord.SendMessageInput{
					ChannelID: taskConfig.ProfitReportChannelID,
					Content:   fmt.Sprintf("string to float64 錯誤: %v", err),
				},
			)
			return
		}
	}

	todayProfit := totalValue - yesterdayTotalValueFloat

	oldMoney := []float64{totalCost, totalValue, totalProfit, todayProfit}
	// 換算幣值
	convertCtx, convertCancel := context.WithTimeout(ctx, externalTimeout)
	newMoney, err := exchange.ConvertExchangeWithContext(convertCtx, oldMoney)
	convertCancel()
	if err != nil {
		logger.Error("換算匯率失敗", "error", err)
		taskErrorReporter.Notify(
			s,
			"calculate_profit:convert_exchange",
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("換算匯率錯誤: %v", err),
			},
		)
		return
	}

	_, err = s.ChannelMessageSendComplex(taskConfig.ProfitReportChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("<@%s> 總成本: %.2f, 目前市場總值: %.2f, 目前損益: %.2f, 今日損益: %.2f  \n 換算台幣總成本: %.2f, 換算台幣目前市場總值: %.2f, 換算台幣目前損益: %.2f, 換算台幣今日損益: %.2f",
			taskConfig.DefaultUserID, totalCost, totalValue, totalProfit, todayProfit, newMoney[0], newMoney[1], newMoney[2], newMoney[3]),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		logger.Error("發送收益報告失敗", "error", err)
		taskErrorReporter.Notify(
			s,
			"calculate_profit:send_report",
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("發送訊息錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("收益報告發送成功", "totalCost", totalCost, "totalValue", totalValue, "totalProfit", totalProfit, "todayProfit", todayProfit)

	// 將今日市場總值存入 Redis
	redisSetCtx, redisSetCancel := context.WithTimeout(ctx, externalTimeout)
	err = redisClient.Set(redisSetCtx, redisKey, totalValue, 0)
	redisSetCancel()
	if err != nil {
		logger.Error("儲存今日市場總值失敗", "error", err)
		taskErrorReporter.Notify(
			s,
			"calculate_profit:redis_set",
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("存今日總值錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("完成收益計算")
}
