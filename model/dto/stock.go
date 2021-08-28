package dto

type Stock struct {
	ID     int64   // 流水號
	UserID string  // 用戶 ID
	Symbol string  // 標的
	Units  float64 // 數量(股)
	Price  float64 // 價格
}
