package discord

import (
	"discordBot/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

// Messenger Discord消息發送接口
type Messenger interface {
	SendMessage(channelID string, content string) error
	SendMessageWithMention(channelID string, userID string, content string) error
}

// messenger 實現Messenger接口
type messenger struct {
	session *discordgo.Session
}

// NewMessenger 創建Messenger實例
func NewMessenger(session *discordgo.Session) Messenger {
	return &messenger{session: session}
}

// SendMessage 發送消息到指定頻道
func (m *messenger) SendMessage(channelID string, content string) error {
	_, err := m.session.ChannelMessageSend(channelID, content)
	if err != nil {
		logger.Error("發送消息失敗", "channelID", channelID, "error", err)
		return err
	}
	return nil
}

// SendMessageWithMention 發送帶有用戶提及的消息
func (m *messenger) SendMessageWithMention(channelID string, userID string, content string) error {
	_, err := m.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		logger.Error("發送提及消息失敗", "channelID", channelID, "userID", userID, "error", err)
		return err
	}
	return nil
}
