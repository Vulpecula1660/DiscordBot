package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"

	"discordBot/handler"
	"discordBot/model/postgresql"
	"discordBot/model/redis"
	"discordBot/pkg/config"
	"discordBot/pkg/logger"
)

func main() {
	// 初始化日誌系統
	logger.Init()
	logger.Info("Discord Bot 啟動中...")

	// 獲取配置
	token := config.GetDiscordToken()
	if token == "" {
		logger.Error("Discord Token 未設置")
		os.Exit(1)
	}

	// 創建 Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("創建 Discord session 失敗", "error", err)
		os.Exit(1)
	}

	// 創建命令路由器
	router := handler.NewCommandRouter()

	// 註冊股票指令
	router.Register("$+", handler.Quote)
	router.Register("$set_stock", handler.SetStock)
	router.Register("$get_stock", handler.GetStock)

	// 註冊 Redis 指令
	router.Register("$setRedis", handler.SetRedis)
	router.Register("$getRedis", handler.GetRedis)
	router.Register("$setList", handler.SetList)
	router.Register("$getList", handler.GetList)
	router.Register("$delListValue", handler.DelListValue)

	// 註冊命令處理器
	dg.AddHandler(router.Handle)

	// 啟動定時任務
	handler.Task(dg)

	// 設置只監聽訊息事件
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// 建立連線
	if err := dg.Open(); err != nil {
		logger.Error("開啟 Discord 連線失敗", "error", err)
		os.Exit(1)
	}

	logger.Info("Bot 已成功啟動，按 CTRL-C 退出")

	// 設置優雅關閉
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// 在 goroutine 中等待關閉信號
	go func() {
		<-sc
		logger.Info("接收到關閉信號，正在優雅關閉...")

		// 創建超時 context
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		// 停止定時任務
		handler.StopTasks()
		logger.Info("定時任務已停止")

		// 關閉資料庫連線
		if err := postgresql.CloseAll(); err != nil {
			logger.Error("關閉資料庫連線失敗", "error", err)
		} else {
			logger.Info("資料庫連線已關閉")
		}

		// 關閉 Redis 連線
		if err := redis.CloseAll(); err != nil {
			logger.Error("關閉 Redis 連線失敗", "error", err)
		} else {
			logger.Info("Redis 連線已關閉")
		}

		// 關閉 Discord session
		if err := dg.Close(); err != nil {
			logger.Error("關閉 Discord session 失敗", "error", err)
		} else {
			logger.Info("Discord session 已關閉")
		}

		// 等待所有操作完成或超時
		select {
		case <-shutdownCtx.Done():
			logger.Warn("優雅關閉超時")
		default:
			logger.Info("優雅關閉完成")
		}

		cancel()
	}()

	// 等待 context 取消
	<-ctx.Done()
	logger.Info("Bot 已退出")
}
