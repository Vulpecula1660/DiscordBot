package stock

import (
	"context"
	"discordBot/model/dao/stock"
	"discordBot/service/exchange"
	"fmt"

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
		s.ChannelMessageSend("872317320729616395", fmt.Sprintf("取資料時錯誤: %v", err))
		return
	}

	var totalProfit, totalCost, totalValue float64

	for _, v := range dbRes {
		value, profit, err := Calculate(
			ctx,
			&CalculateInput{
				Symbol: v.Symbol,
				Units:  v.Units,
				Price:  v.Price,
			})
		if err != nil {
			s.ChannelMessageSend("872317320729616395", fmt.Sprintf("計算損益時錯誤: %v", err))
			return
		}
		cost := v.Units * v.Price

		totalProfit = totalProfit + profit
		totalCost = totalCost + cost
		totalValue = totalValue + value
	}

	oldMoney := []float64{totalCost, totalValue, totalProfit}
	newMoney, err := exchange.ConvertExchange(oldMoney)
	if err != nil {
		s.ChannelMessageSend("872317320729616395", fmt.Sprintf("換算匯率錯誤: %v", err))
		return
	}

	_, err = s.ChannelMessageSendComplex("872317320729616395", &discordgo.MessageSend{
		Content: fmt.Sprintf("<@512265930735222795> 總成本: %.2f, 目前市場總值: %.2f, 目前損益: %.2f \n 換算台幣總成本: %.2f, 換算台幣目前市場總值: %.2f, 換算台幣目前損益: %.2f", totalCost, totalValue, totalProfit, newMoney[0], newMoney[1], newMoney[2]),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		s.ChannelMessageSend("872317320729616395", fmt.Sprintf("發送訊息錯誤: %v", err))
		return
	}
}
