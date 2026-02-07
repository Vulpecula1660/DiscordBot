package handler

import (
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler 命令處理函數類型
type CommandHandler func(*discordgo.Session, *discordgo.MessageCreate)

// commandEntry 命令路由條目
type commandEntry struct {
	prefix  string
	handler CommandHandler
}

// CommandRouter 命令路由器
type CommandRouter struct {
	commands []commandEntry
}

// NewCommandRouter 創建命令路由器
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{}
}

// Register 註冊命令處理器
func (r *CommandRouter) Register(prefix string, handler CommandHandler) {
	r.commands = append(r.commands, commandEntry{prefix: prefix, handler: handler})
	// 依前綴長度降序排列，確保最長前綴優先匹配
	sort.Slice(r.commands, func(i, j int) bool {
		return len(r.commands[i].prefix) > len(r.commands[j].prefix)
	})
}

// Handle 處理消息事件
func (r *CommandRouter) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 忽略機器人自己的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 依最長前綴優先順序匹配
	for _, entry := range r.commands {
		if strings.HasPrefix(m.Content, entry.prefix) {
			entry.handler(s, m)
			return
		}
	}
}

// GetRegisteredCommands 獲取已註冊的命令列表（用於調試）
func (r *CommandRouter) GetRegisteredCommands() []string {
	commands := make([]string, 0, len(r.commands))
	for _, entry := range r.commands {
		commands = append(commands, entry.prefix)
	}
	return commands
}
