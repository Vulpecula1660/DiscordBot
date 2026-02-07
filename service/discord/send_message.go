package discord

import (
	"github.com/bwmarrin/discordgo"
)

type SendMessageInput struct {
	ChannelID string
	Content   string
}

// SendMessage : 發送消息到指定頻道
func SendMessage(s Session, input *SendMessageInput) error {
	_, err := s.ChannelMessageSend(input.ChannelID, input.Content)
	if err != nil {
		return err
	}

	return nil
}

// Ensure *discordgo.Session implements Session interface
var _ Session = (*discordgo.Session)(nil)
