package config

import (
	"os"
	"strconv"
)

// GetDiscordToken 獲取Discord Token
func GetDiscordToken() string {
	return os.Getenv("DCToken")
}

// TaskConfig 定時任務相關配置
type TaskConfig struct {
	// 加密貨幣價格更新頻道
	CryptoPriceChannelID string
	// 股票觀察清單頻道
	WatchListChannelID string
	// 收益報告頻道
	ProfitReportChannelID string
	// 默認用戶ID（用於收益報告）
	DefaultUserID string
	// CheckChange 任務併發上限
	CheckChangeMaxConcurrency int
	// CalculateProfit 任務併發上限
	CalculateProfitMaxConcurrency int
	// CheckChange 任務總超時（秒）
	CheckChangeTimeoutSeconds int
	// CalculateProfit 任務總超時（秒）
	CalculateProfitTimeoutSeconds int
	// 單次外部呼叫超時（秒）
	ExternalCallTimeoutSeconds int
	// 任務錯誤通知節流間隔（秒）
	ErrorNotifyCooldownSeconds int
}

// GetTaskConfig 獲取定時任務配置
func GetTaskConfig() *TaskConfig {
	return &TaskConfig{
		CryptoPriceChannelID:          getEnv("CRYPTO_PRICE_CHANNEL_ID", ""),
		WatchListChannelID:            getEnv("WATCH_LIST_CHANNEL_ID", ""),
		ProfitReportChannelID:         getEnv("PROFIT_REPORT_CHANNEL_ID", ""),
		DefaultUserID:                 getEnv("DEFAULT_USER_ID", ""),
		CheckChangeMaxConcurrency:     getEnvInt("TASK_CHECK_CHANGE_MAX_CONCURRENCY", 5),
		CalculateProfitMaxConcurrency: getEnvInt("TASK_CALCULATE_PROFIT_MAX_CONCURRENCY", 5),
		CheckChangeTimeoutSeconds:     getEnvInt("TASK_CHECK_CHANGE_TIMEOUT_SECONDS", 120),
		CalculateProfitTimeoutSeconds: getEnvInt("TASK_CALCULATE_PROFIT_TIMEOUT_SECONDS", 180),
		ExternalCallTimeoutSeconds:    getEnvInt("TASK_EXTERNAL_CALL_TIMEOUT_SECONDS", 15),
		ErrorNotifyCooldownSeconds:    getEnvInt("TASK_ERROR_NOTIFY_COOLDOWN_SECONDS", 60),
	}
}

// Helper functions
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
