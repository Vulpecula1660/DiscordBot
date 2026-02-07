package stock

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordBot/model/dao/stock"
	"discordBot/model/dto"
)

// SetStock : 將股票新增到 DB
func SetStock(ctx context.Context, m *discordgo.MessageCreate) error {
	// example : $set_stock TSLA units price
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 4 {
		return fmt.Errorf("參數錯誤，格式: $set_stock <symbol> <units> <price>")
	}

	symbol := strSlice[1]
	unitsStr := strSlice[2]
	priceStr := strSlice[3]

	units, err := strconv.ParseFloat(unitsStr, 64)
	if err != nil {
		return fmt.Errorf("無效的數量: %w", err)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("無效的價格: %w", err)
	}

	if err := stock.Ins(
		ctx,
		nil,
		&dto.Stock{
			UserID: m.Author.ID,
			Symbol: symbol,
			Units:  units,
			Price:  price,
		},
	); err != nil {
		return err
	}

	return nil
}

// GetStock : DB 取得股票
func GetStock(ctx context.Context, m *discordgo.MessageCreate) ([]*dto.Stock, error) {
	// example : $get_stock TSLA
	strSlice := strings.Fields(m.Content)

	if len(strSlice) != 2 {
		return nil, fmt.Errorf("參數錯誤，格式: $get_stock <symbol>")
	}

	symbol := strSlice[1]

	res, err := stock.Get(
		ctx,
		&stock.GetInput{
			Symbol: symbol,
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
