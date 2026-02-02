package stock

import (
	"context"
	"fmt"
	"os"
	"strings"

	"discordBot/model/dto"
	"discordBot/model/postgresql"
)

// GetInput :
type GetInput struct {
	UserID string
	Symbol string
}

// Get : 取得 d9fdq7n9q3delq.stock
func Get(ctx context.Context, input *GetInput) (ret []*dto.Stock, err error) {
	dbS, err := postgresql.GetConn(os.Getenv("DATABASE_Name"))
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var params []interface{}
	var wheres []string
	paramIndex := 1

	// UserID
	if input.UserID != "" {
		wheres = append(wheres, fmt.Sprintf(" user_id = $%d ", paramIndex))
		params = append(params, input.UserID)
		paramIndex++
	}

	// Symbol
	if input.Symbol != "" {
		wheres = append(wheres, fmt.Sprintf(" symbol = $%d ", paramIndex))
		params = append(params, input.Symbol)
		paramIndex++
	}

	// 沒有條件時回傳錯誤
	if len(wheres) == 0 {
		return nil, fmt.Errorf("sql 語法錯誤")
	}

	sql := `SELECT id, user_id, symbol, units, price FROM stock WHERE` + strings.Join(wheres, " AND ")

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
