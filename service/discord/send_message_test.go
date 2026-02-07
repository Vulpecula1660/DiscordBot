package discord

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// MockSession Session 的 mock 實現
type MockSession struct {
	Messages []MockMessage
	Err      error
}

// MockMessage 記錄發送的消息
type MockMessage struct {
	ChannelID string
	Content   string
}

// ChannelMessageSend 實現 Session 接口
func (m *MockSession) ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	m.Messages = append(m.Messages, MockMessage{
		ChannelID: channelID,
		Content:   content,
	})

	return &discordgo.Message{}, nil
}

// NewMockSession 創建一個新的 mock session
func NewMockSession() *MockSession {
	return &MockSession{
		Messages: make([]MockMessage, 0),
	}
}

func Test_SendMessage(t *testing.T) {
	tests := []struct {
		name          string
		input         *SendMessageInput
		mockErr       error
		wantErr       bool
		wantMsgCount  int
		wantChannelID string
		wantContent   string
	}{
		{
			name: "successful send",
			input: &SendMessageInput{
				ChannelID: "123456789",
				Content:   "Test message",
			},
			wantErr:       false,
			wantMsgCount:  1,
			wantChannelID: "123456789",
			wantContent:   "Test message",
		},
		{
			name: "empty content",
			input: &SendMessageInput{
				ChannelID: "123456789",
				Content:   "",
			},
			wantErr:       false,
			wantMsgCount:  1,
			wantChannelID: "123456789",
			wantContent:   "",
		},
		{
			name: "empty channel ID",
			input: &SendMessageInput{
				ChannelID: "",
				Content:   "Test message",
			},
			wantErr:       false,
			wantMsgCount:  1,
			wantChannelID: "",
			wantContent:   "Test message",
		},
		{
			name: "send error",
			input: &SendMessageInput{
				ChannelID: "123456789",
				Content:   "Test message",
			},
			mockErr:      errors.New("connection failed"),
			wantErr:      true,
			wantMsgCount: 0,
		},
		{
			name: "long content",
			input: &SendMessageInput{
				ChannelID: "123456789",
				Content:   "This is a very long message with special characters: !@#$%^&*()_+-=[]{}|;':\",./<>?",
			},
			wantErr:       false,
			wantMsgCount:  1,
			wantChannelID: "123456789",
			wantContent:   "This is a very long message with special characters: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockSession()
			mock.Err = tt.mockErr

			err := SendMessage(mock, tt.input)

			// 驗證錯誤
			if tt.wantErr {
				if err == nil {
					t.Errorf("SendMessage() error = nil, wantErr = true")
					return
				}
				if tt.mockErr != nil && err.Error() != tt.mockErr.Error() {
					t.Errorf("SendMessage() error = %v, want %v", err.Error(), tt.mockErr.Error())
				}
			} else {
				if err != nil {
					t.Errorf("SendMessage() unexpected error = %v", err)
					return
				}
			}

			// 驗證發送的消息數量
			if len(mock.Messages) != tt.wantMsgCount {
				t.Errorf("SendMessage() sent %d messages, want %d", len(mock.Messages), tt.wantMsgCount)
				return
			}

			// 驗證消息內容
			if tt.wantMsgCount > 0 {
				msg := mock.Messages[0]
				if msg.ChannelID != tt.wantChannelID {
					t.Errorf("SendMessage() channelID = %v, want %v", msg.ChannelID, tt.wantChannelID)
				}
				if msg.Content != tt.wantContent {
					t.Errorf("SendMessage() content = %v, want %v", msg.Content, tt.wantContent)
				}
			}
		})
	}
}

func Test_SendMessage_Multiple(t *testing.T) {
	mock := NewMockSession()

	messages := []*SendMessageInput{
		{ChannelID: "111", Content: "First message"},
		{ChannelID: "222", Content: "Second message"},
		{ChannelID: "333", Content: "Third message"},
	}

	for _, msg := range messages {
		if err := SendMessage(mock, msg); err != nil {
			t.Errorf("SendMessage() unexpected error = %v", err)
		}
	}

	if len(mock.Messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(mock.Messages))
	}

	for i, msg := range mock.Messages {
		if msg.ChannelID != messages[i].ChannelID {
			t.Errorf("Message %d: channelID = %v, want %v", i, msg.ChannelID, messages[i].ChannelID)
		}
		if msg.Content != messages[i].Content {
			t.Errorf("Message %d: content = %v, want %v", i, msg.Content, messages[i].Content)
		}
	}
}
