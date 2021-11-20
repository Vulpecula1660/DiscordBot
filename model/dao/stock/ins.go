package stock

import (
	"context"
	dbSQL "database/sql"
	"fmt"
	"os"

	"discordBot/model/dto"
	"discordBot/model/postgresql"
)

// Ins : 新增股票 ins d9fdq7n9q3delq.stock
// Transaction 為選填
func Ins(ctx context.Context, tx *dbSQL.Tx, input *dto.Stock) (err error) {
	if input == nil {
		return fmt.Errorf("參數錯誤")
	}

	var dbM *dbSQL.DB

	if tx == nil {
		dbM = postgresql.GetConn(os.Getenv("DATABASE_Name"))
	}

	sql := " INSERT INTO stock ("
	sql += "    user_id,"
	sql += "    symbol,"
	sql += "    units,"
	sql += "    price "
	sql += " )"
	sql += " VALUES "
	sql += " ( $1, $2, $3, $4)"

	var params []interface{}

	params = append(params, input.UserID)
	params = append(params, input.Symbol)
	params = append(params, input.Units)
	params = append(params, input.Price)

	// 執行sql

	if tx == nil {
		_, err = dbM.ExecContext(ctx, sql, params...)
	} else {
		_, err = tx.ExecContext(ctx, sql, params...)
	}
	if err != nil {
		return fmt.Errorf("ins錯誤 error: %v, sql: %v, params: %v ", err, sql, params)
	}

	return err
}
