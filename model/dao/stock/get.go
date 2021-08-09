package stock

import (
	"context"
	"discordBot/model/dto"
	"discordBot/model/postgresql"
	"fmt"
	"os"
	"strings"
)

// GetInput :
type GetInput struct {
	UserID string
}

// Get : 取得 d9fdq7n9q3delq.stock
func Get(ctx context.Context, input *GetInput) (ret []*dto.Stock, err error) {
	dbS := postgresql.GetConn(os.Getenv("DATABASE_Name"))

	sql := " SELECT "
	sql += "    id, "
	sql += "    user_id, "
	sql += "    symbol, "
	sql += "    units, "
	sql += "    price "

	sql += " FROM stock "
	sql += " WHERE "

	var params []interface{}
	var wheres []string

	// UserID
	if input.UserID != "" {
		wheres = append(wheres, " user_id = $1 ")
		params = append(params, input.UserID)
	}

	// 沒有條件時回傳錯誤
	if len(wheres) == 0 {
		return nil, fmt.Errorf("sql 語法錯誤")
	}

	sql += strings.Join(wheres, " AND ")

	rows, err := dbS.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, fmt.Errorf("select 錯誤: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		data := &dto.Stock{}
		if err := rows.Scan(
			&data.ID,
			&data.UserID,
			&data.Symbol,
			&data.Units,
			&data.Price,
		); err != nil {
			return nil, fmt.Errorf("scan 錯誤: %v", err)
		}
		ret = append(ret, data)
	}

	return ret, err
}
