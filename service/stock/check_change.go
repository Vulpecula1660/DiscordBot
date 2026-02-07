package stock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"discordBot/model/redis"
	"discordBot/pkg/config"
	"discordBot/pkg/logger"
	"discordBot/service/discord"

	"github.com/bwmarrin/discordgo"
)

// RedisClient Redis 客戶端接口類型
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}

// CheckChange : 檢查漲跌幅
func CheckChange(s *discordgo.Session) {
	CheckChangeWithDeps(s, redisDeps{})
}

// redisDeps 封裝 Redis 依賴
type redisDeps struct{}

func (d redisDeps) Get(ctx context.Context, key string) (string, error) {
	return redis.Get(ctx, key)
}

func (d redisDeps) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return redis.Set(ctx, key, value, expiration)
}

func (d redisDeps) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return redis.LRange(ctx, key, start, stop)
}

// CheckChangeWithDeps 使用指定依賴檢查漲跌幅（用於測試）
func CheckChangeWithDeps(s *discordgo.Session, redisClient RedisClient) {
	taskConfig := config.GetTaskConfig()
	taskErrorReporter.SetCooldown(durationFromSeconds(taskConfig.ErrorNotifyCooldownSeconds, time.Minute))

	runTimeout := durationFromSeconds(taskConfig.CheckChangeTimeoutSeconds, 2*time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()

	externalTimeout := durationFromSeconds(taskConfig.ExternalCallTimeoutSeconds, 15*time.Second)
	maxConcurrency := normalizeConcurrency(taskConfig.CheckChangeMaxConcurrency, 5)

	logger.Info("開始檢查股票漲跌幅", "maxConcurrency", maxConcurrency, "timeout", runTimeout.String())

	// Redis 取出資料
	fetchCtx, fetchCancel := context.WithTimeout(ctx, externalTimeout)
	watchList, err := redisClient.LRange(
		fetchCtx,
		"watch_list",
		0,
		-1,
	)
	fetchCancel()
	if err != nil {
		logger.Error("取得觀察列表失敗", "error", err)
		taskErrorReporter.Notify(
			s,
			"check_change:watch_list:lrange",
			&discord.SendMessageInput{
				ChannelID: taskConfig.WatchListChannelID,
				Content:   fmt.Sprintf("取得列表時錯誤: %v", err),
			},
		)
		return
	}

	logger.Info("取得觀察列表", "count", len(watchList))
	if len(watchList) == 0 {
		logger.Info("觀察列表為空，略過漲跌幅檢查")
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

loop:
	for _, v := range watchList {
		if err := ctx.Err(); err != nil {
			logger.Warn("檢查任務超時，停止派發剩餘標的", "error", err)
			break
		}

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			logger.Warn("檢查任務超時，停止派發剩餘標的", "error", ctx.Err())
			break loop
		}

		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			// 先看是否已通知過
			redisGetCtx, redisGetCancel := context.WithTimeout(ctx, externalTimeout)
			redisRes, err := redisClient.Get(redisGetCtx, "watch_list:"+symbol)
			redisGetCancel()
			if err != nil {
				logger.Error("取得通知紀錄失敗", "symbol", symbol, "error", err)
				taskErrorReporter.Notify(
					s,
					"check_change:watch_list:get_record",
					&discord.SendMessageInput{
						ChannelID: taskConfig.WatchListChannelID,
						Content:   fmt.Sprintf("取得紀錄時錯誤: %v", err),
					},
				)
				return
			}

			if redisRes != "" {
				return
			}

			changeCtx, changeCancel := context.WithTimeout(ctx, externalTimeout)
			change, err := GetChange(changeCtx, symbol)
			changeCancel()
			if err != nil {
				logger.Error("取得漲跌幅失敗", "symbol", symbol, "error", err)
				taskErrorReporter.Notify(
					s,
					"check_change:watch_list:get_change",
					&discord.SendMessageInput{
						ChannelID: taskConfig.WatchListChannelID,
						Content:   fmt.Sprintf("取得漲跌幅時錯誤: %v", err),
					},
				)
				return
			}

			if change > 3 || change < -3 {
				logger.Warn("股票漲跌幅超過閾值", "symbol", symbol, "change", change)
				_, err = s.ChannelMessageSendComplex(taskConfig.WatchListChannelID, &discordgo.MessageSend{
					Content: fmt.Sprintf("<@%s> 警告: %s 今日漲跌幅為 %.2f %%", taskConfig.DefaultUserID, symbol, change),
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
					},
				})
				if err != nil {
					logger.Error("發送警告訊息失敗", "symbol", symbol, "error", err)
					taskErrorReporter.Notify(
						s,
						"check_change:watch_list:send_alert",
						&discord.SendMessageInput{
							ChannelID: taskConfig.WatchListChannelID,
							Content:   fmt.Sprintf("發送訊息時錯誤: %v", err),
						},
					)
					return
				}

				// 寫入紀錄已通知
				setCtx, setCancel := context.WithTimeout(ctx, externalTimeout)
				err = redisClient.Set(setCtx, "watch_list:"+symbol, "true", time.Hour*8)
				setCancel()
				if err != nil {
					logger.Error("寫入通知紀錄失敗", "symbol", symbol, "error", err)
					taskErrorReporter.Notify(
						s,
						"check_change:watch_list:set_record",
						&discord.SendMessageInput{
							ChannelID: taskConfig.WatchListChannelID,
							Content:   fmt.Sprintf("寫入紀錄時錯誤: %v", err),
						},
					)
					return
				}
			}
		}(v)
	}

	wg.Wait()
	if err := ctx.Err(); err != nil && err != context.Canceled {
		logger.Warn("股票漲跌幅檢查任務逾時", "error", err)
		taskErrorReporter.Notify(
			s,
			"check_change:timeout",
			&discord.SendMessageInput{
				ChannelID: taskConfig.WatchListChannelID,
				Content:   fmt.Sprintf("股票漲跌幅檢查任務逾時: %v", err),
			},
		)
	}

	logger.Info("完成股票漲跌幅檢查")
}
