package stock

import (
	"context"
	"discordBot/model/dao/stock"
	"discordBot/model/dto"
	"discordBot/model/redis"
	"discordBot/service/discord"
	"discordBot/service/exchange"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// CalculateProfit : 計算損益
func CalculateProfit(s *discordgo.Session) {
	ctx := context.Background()

	// 先取資料
	dbRes, err := stock.Get(
		ctx,
		&stock.GetInput{
			UserID: "512265930735222795",
		},
	)
	if err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: "872317320729616395",
				Content:   fmt.Sprintf("取資料時錯誤: %v", err),
			},
		)
		return
	}

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
				discord.SendMessage(
					s,
					&discord.SendMessageInput{
						ChannelID: "872317320729616395",
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
	yesterdayTotalValue, err := redis.Get(ctx, "872317320729616395_"+"totalValue")
	if err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: "872317320729616395",
				Content:   fmt.Sprintf("取昨日市場總值錯誤: %v", err),
			},
		)
		return
	}

	// string to float64
	yesterdayTotalValueFloat, err := strconv.ParseFloat(yesterdayTotalValue, 64)
	if err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: "872317320729616395",
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
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: "872317320729616395",
				Content:   fmt.Sprintf("換算匯率錯誤: %v", err),
			},
		)
		return
	}

	_, err = s.ChannelMessageSendComplex("872317320729616395", &discordgo.MessageSend{
		Content: fmt.Sprintf("<@512265930735222795> 總成本: %.2f, 目前市場總值: %.2f, 目前損益: %.2f, 今日損益: %.2f  \n 換算台幣總成本: %.2f, 換算台幣目前市場總值: %.2f, 換算台幣目前損益: %.2f, 換算台幣今日損益: %.2f", totalCost, totalValue, totalProfit, todayProfit, newMoney[0], newMoney[1], newMoney[2], newMoney[3]),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: "872317320729616395",
				Content:   fmt.Sprintf("發送訊息錯誤: %v", err),
			},
		)
		return
	}

	// 將今日市場總值存入 Redis
	err = redis.Set(ctx, "872317320729616395_"+"totalValue", totalValue, time.Hour*0)
	if err != nil {
		discord.SendMessage(
			s,
			&discord.SendMessageInput{
				ChannelID: "872317320729616395",
				Content:   fmt.Sprintf("存今日總值錯誤: %v", err),
			},
		)
		return
	}
}
