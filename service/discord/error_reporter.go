package discord

import (
	"fmt"
	"sync"
	"time"

	"discordBot/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

// ErrorReporter 針對高頻錯誤通知做節流與去重
type ErrorReporter struct {
	cooldown time.Duration

	mu       sync.Mutex
	lastSent map[string]time.Time
}

// NewErrorReporter 建立錯誤通知器
func NewErrorReporter(cooldown time.Duration) *ErrorReporter {
	if cooldown <= 0 {
		cooldown = time.Minute
	}

	return &ErrorReporter{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

// SetCooldown 更新節流時間
func (r *ErrorReporter) SetCooldown(cooldown time.Duration) {
	if r == nil {
		return
	}
	if cooldown <= 0 {
		cooldown = time.Minute
	}

	r.mu.Lock()
	r.cooldown = cooldown
	r.mu.Unlock()
}

// Notify 發送節流後的錯誤通知
func (r *ErrorReporter) Notify(s *discordgo.Session, key string, input *SendMessageInput) {
	if r == nil || s == nil || input == nil || input.ChannelID == "" {
		return
	}

	if key == "" {
		key = fmt.Sprintf("%s:%s", input.ChannelID, input.Content)
	}

	if !r.shouldSend(key) {
		logger.Warn("錯誤通知已節流", "key", key, "cooldown", r.cooldown.String())
		return
	}

	if err := SendMessage(s, input); err != nil {
		logger.Error("發送錯誤通知失敗", "key", key, "channelID", input.ChannelID, "error", err)
	}
}

func (r *ErrorReporter) shouldSend(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	last, ok := r.lastSent[key]
	if ok && now.Sub(last) < r.cooldown {
		return false
	}

	r.lastSent[key] = now
	return true
}
