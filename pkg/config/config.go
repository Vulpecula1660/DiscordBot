package config

import (
	"os"
	"strconv"
)

// DiscordConfig 包含Discord相關配置
type DiscordConfig struct {
	Token string
}

// GetDiscordToken 獲取Discord Token
func GetDiscordToken() string {
	return os.Getenv("DCToken")
}

// ChannelConfig 包含各種頻道ID配置
type ChannelConfig struct {
	// StockChannels 股票相關頻道
	StockChannels StockChannels
	// WatchListChannels 觀察清單頻道
	WatchListChannels WatchListChannels
	// CryptoChannels 加密貨幣相關頻道
	CryptoChannels CryptoChannels
}

// StockChannels 股票頻道配置
type StockChannels struct {
	ProfitReport string
	PriceCheck   string
}

// WatchListChannels 觀察清單頻道配置
type WatchListChannels struct {
	WatchList string
}

// CryptoChannels 加密貨幣頻道配置
type CryptoChannels struct {
	PriceUpdate string
}

// GetChannelConfig 獲取頻道配置
func GetChannelConfig() *ChannelConfig {
	return &ChannelConfig{
		StockChannels: StockChannels{
			ProfitReport: getEnv("CHANNEL_PROFIT_REPORT", ""),
			PriceCheck:   getEnv("CHANNEL_PRICE_CHECK", ""),
		},
		WatchListChannels: WatchListChannels{
			WatchList: getEnv("CHANNEL_WATCH_LIST", ""),
		},
		CryptoChannels: CryptoChannels{
			PriceUpdate: getEnv("CHANNEL_CRYPTO_UPDATE", ""),
		},
	}
}

// UserConfig 用戶相關配置
type UserConfig struct {
	// DefaultUserID 默認用戶ID（用於股票收益報告）
	DefaultUserID string
}

// GetUserConfig 獲取用戶配置
func GetUserConfig() *UserConfig {
	return &UserConfig{
		DefaultUserID: getEnv("DEFAULT_USER_ID", ""),
	}
}

// DatabaseConfig 數據庫配置
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// GetDatabaseConfig 獲取數據庫配置
func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     os.Getenv("DATABASE_Host"),
		Port:     os.Getenv("DATABASE_Port"),
		Name:     os.Getenv("DATABASE_Name"),
		User:     os.Getenv("DATABASE_User"),
		Password: os.Getenv("DATABASE_Password"),
	}
}

// RedisConfig Redis配置
type RedisConfig struct {
	URL string
}

// GetRedisConfig 獲取Redis配置
func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		URL: os.Getenv("REDIS_URL"),
	}
}

// APIConfig 外部API配置
type APIConfig struct {
	FinnhubAPIKey string
}

// GetAPIConfig 獲取API配置
func GetAPIConfig() *APIConfig {
	return &APIConfig{
		FinnhubAPIKey: os.Getenv("APIKey"),
	}
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
}

// GetTaskConfig 獲取定時任務配置
func GetTaskConfig() *TaskConfig {
	return &TaskConfig{
		CryptoPriceChannelID:  getEnv("CRYPTO_PRICE_CHANNEL_ID", "1032641300077490266"),
		WatchListChannelID:    getEnv("WATCH_LIST_CHANNEL_ID", "960897897166176266"),
		ProfitReportChannelID: getEnv("PROFIT_REPORT_CHANNEL_ID", "872317320729616395"),
		DefaultUserID:         getEnv("DEFAULT_USER_ID", "512265930735222795"),
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

func getEnvFloat64(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
			return floatVal
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultVal
}
