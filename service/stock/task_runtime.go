package stock

import (
	"time"

	"discordBot/service/discord"
)

var taskErrorReporter = discord.NewErrorReporter(time.Minute)

func durationFromSeconds(seconds int, defaultDuration time.Duration) time.Duration {
	if seconds <= 0 {
		return defaultDuration
	}

	return time.Duration(seconds) * time.Second
}

func normalizeConcurrency(value int, fallback int) int {
	if value <= 0 {
		return fallback
	}

	return value
}
