package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

// Init 初始化日誌系統
func Init() {
	once.Do(func() {
		// 根據環境選擇日誌格式
		var handler slog.Handler
		if os.Getenv("ENV") == "production" {
			// 生產環境使用JSON格式
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     getLogLevel(),
				AddSource: true,
			})
		} else {
			// 開發環境使用文本格式
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     getLogLevel(),
				AddSource: false,
			})
		}
		defaultLogger = slog.New(handler)
		slog.SetDefault(defaultLogger)
	})
}

// getLogLevel 從環境變量獲取日誌級別
func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Debug 輸出DEBUG級別日誌
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Info 輸出INFO級別日誌
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Warn 輸出WARN級別日誌
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Error 輸出ERROR級別日誌
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// With 創建帶屬性的Logger
func With(args ...any) *slog.Logger {
	return slog.With(args...)
}
