package stock

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"discordBot/model/dao/stock"
	"discordBot/model/dto"
	"discordBot/model/redis"
	"discordBot/pkg/config"
	"discordBot/pkg/logger"
	"discordBot/service/discord"
	"discordBot/service/exchange"

	"github.com/bwmarrin/discordgo"
)

// CalculateProfit : 計算損益
func CalculateProfit(s *discordgo.Session) {
	ctx := context.Background()
	taskConfig := config.GetTaskConfig()

	logger.Info("開始計算收益")

	// 先取資料
	dbRes, err := stock.Get(
		ctx,
		&stock.GetInput{
			UserID: taskConfig.DefaultUserID,
		},
	)
	if err != nil {
		logger.Error("取得股票資料失敗", "error", err)
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("取資料時錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("取得股票資料", "count", len(dbRes))

	var totalProfit, totalCost, totalValue float64

	var wg sync.WaitGroup
	wg.Add(len(dbRes))

	var mu sync.Mutex

	for _, v := range dbRes {
		go func(stock *dto.Stock) {
			defer wg.Done()

			value, profit, err := Calculate(
				ctx,
				&CalculateInput{
					Symbol: stock.Symbol,
					Units:  stock.Units,
					Price:  stock.Price,
				})
			if err != nil {
				logger.Error("計算損益失敗", "symbol", stock.Symbol, "error", err)
				discord.SendMessage(
					s,
					&discord.SendMessageInput{
						ChannelID: taskConfig.ProfitReportChannelID,
						Content:   fmt.Sprintf("計算損益時錯誤: %v", err),
					},
				)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			cost := stock.Units * stock.Price

			totalProfit = totalProfit + profit
			totalCost = totalCost + cost
			totalValue = totalValue + value
		}(v)
	}

	wg.Wait()

	// 取昨日市場總值
	redisKey := taskConfig.ProfitReportChannelID + "_" + "totalValue"
	yesterdayTotalValue, err := redis.Get(ctx, redisKey)
	if err != nil {
		logger.Error("取得昨日市場總值失敗", "error", err)
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("取昨日市場總值錯誤: %v", err),
			},
		)
		return
	}

	// string to float64
	yesterdayTotalValueFloat, err := strconv.ParseFloat(yesterdayTotalValue, 64)
	if err != nil {
		logger.Error("轉換昨日市場總值失敗", "value", yesterdayTotalValue, "error", err)
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("string to float64 錯誤: %v", err),
			},
		)
		return
	}

	todayProfit := totalValue - yesterdayTotalValueFloat

	oldMoney := []float64{totalCost, totalValue, totalProfit, todayProfit}
	// 換算幣值
	newMoney, err := exchange.ConvertExchange(oldMoney)
	if err != nil {
		logger.Error("換算匯率失敗", "error", err)
		discord.SendMessage(
			s,
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
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("發送訊息錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("收益報告發送成功", "totalCost", totalCost, "totalValue", totalValue, "totalProfit", totalProfit, "todayProfit", todayProfit)

	// 將今日市場總值存入 Redis
	err = redis.Set(ctx, redisKey, totalValue, time.Hour*0)
	if err != nil {
		logger.Error("儲存今日市場總值失敗", "error", err)
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: taskConfig.ProfitReportChannelID,
				Content:   fmt.Sprintf("存今日總值錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("完成收益計算")
}
