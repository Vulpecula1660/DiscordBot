package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// 預編譯正則表達式
var (
	// 允許的 Redis key 字元：字母、數字、底線、冒號、連字號
	validRedisKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9_:\-]+$`)

	// 股票代號：1-10 個大寫字母或數字
	validStockSymbolPattern = regexp.MustCompile(`^[A-Z0-9]{1,10}$`)
)

// ValidateRedisKey 驗證 Redis key 是否合法
func ValidateRedisKey(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("key cannot be empty")
	}

	if len(key) > 256 {
		return fmt.Errorf("key length must not exceed 256 characters")
	}

	// 禁止特殊字元
	if strings.ContainsAny(key, "\n\r\t ") {
		return fmt.Errorf("key contains invalid characters (whitespace)")
	}

	if !validRedisKeyPattern.MatchString(key) {
		return fmt.Errorf("key contains invalid characters, only alphanumeric, underscore, colon, and hyphen are allowed")
	}

	return nil
}

// ValidateRedisValue 驗證 Redis value 是否合法
func ValidateRedisValue(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("value cannot be empty")
	}

	if len(value) > 1024 {
		return fmt.Errorf("value length must not exceed 1024 characters")
	}

	return nil
}

// ValidateStockSymbol 驗證股票代號是否合法
func ValidateStockSymbol(symbol string) error {
	if len(symbol) == 0 {
		return fmt.Errorf("stock symbol cannot be empty")
	}

	upperSymbol := strings.ToUpper(symbol)
	if !validStockSymbolPattern.MatchString(upperSymbol) {
		return fmt.Errorf("invalid stock symbol format")
	}

	return nil
}

// SanitizeInput 清理使用者輸入，移除潛在危險字元
func SanitizeInput(input string) string {
	// 移除控制字元
	sanitized := strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, input)

	// 移除前後空白
	return strings.TrimSpace(sanitized)
}
