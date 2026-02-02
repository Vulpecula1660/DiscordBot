package handler

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler 命令處理函數類型
type CommandHandler func(*discordgo.Session, *discordgo.MessageCreate)

// CommandRouter 命令路由器
type CommandRouter struct {
	commands map[string]CommandHandler
}

// NewCommandRouter 創建命令路由器
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commands: make(map[string]CommandHandler),
	}
}

// Register 註冊命令處理器
func (r *CommandRouter) Register(prefix string, handler CommandHandler) {
	r.commands[prefix] = handler
}

// Handle 處理消息事件
func (r *CommandRouter) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 忽略機器人自己的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 遍歷所有註冊的命令前綴
	for prefix, handler := range r.commands {
		if strings.HasPrefix(m.Content, prefix) {
			handler(s, m)
			return
		}
	}
}

// GetRegisteredCommands 獲取已註冊的命令列表（用於調試）
func (r *CommandRouter) GetRegisteredCommands() []string {
	commands := make([]string, 0, len(r.commands))
	for prefix := range r.commands {
		commands = append(commands, prefix)
	}
	return commands
}
