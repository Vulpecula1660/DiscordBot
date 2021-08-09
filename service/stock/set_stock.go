package stock

import (
	"context"
	"discordBot/model/dao/stock"
	"discordBot/model/dto"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SetStock(m *discordgo.MessageCreate) error {
	strSlice := strings.Split(m.Content, " ")

	symbol := strSlice[1]
	unitsStr := strSlice[2]
	priceStr := strSlice[3]

	units, _ := strconv.ParseFloat(unitsStr, 64)
	price, _ := strconv.ParseFloat(priceStr, 64)

	err := stock.Ins(
		context.Background(),
		nil,
		&dto.Stock{
			UserID: m.Author.ID,
			Symbol: symbol,
			Units:  units,
			Price:  price,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func GetStock(m *discordgo.MessageCreate) ([]*dto.Stock, error) {
	strSlice := strings.Split(m.Content, " ")

	userID := strSlice[1]

	res, err := stock.Get(
		context.Background(),
		&stock.GetInput{
			UserID: userID,
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
