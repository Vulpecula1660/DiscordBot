package discord

import "github.com/bwmarrin/discordgo"

type SendMessageInput struct {
	ChannelID string
	Content   string
}

func SendMessage(s *discordgo.Session, input *SendMessageInput) error {
	_, err := s.ChannelMessageSend(input.ChannelID, input.Content)
	if err != nil {
		return err
	}

	return nil
}
