package discord

import "github.com/bwmarrin/discordgo"

// Session Discord session 接口（用於測試）
type Session interface {
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}
