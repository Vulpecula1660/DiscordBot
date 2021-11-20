package discord

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func Test_SendMessage(t *testing.T) {
	token := os.Getenv("DCToken")

	// creates a new Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		t.Error("error creating Discord session,", err)
		return
	}

	type args struct {
		s     *discordgo.Session
		input *SendMessageInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test1",
			args: args{
				s: dg,
				input: &SendMessageInput{
					ChannelID: "872317320729616395",
					Content:   "Test1 測試發送",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMessage(tt.args.s, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
