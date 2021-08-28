package main

import (
	"discordBot/handler"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	token := os.Getenv("DCToken")

	// creates a new Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	// 股票指令
	dg.AddHandler(handler.Quote)

	// Redis 指令
	dg.AddHandler(handler.SetRedis)
	dg.AddHandler(handler.GetRedis)
	dg.AddHandler(handler.SetList)
	dg.AddHandler(handler.GetList)
	dg.AddHandler(handler.DelListValue)

	// DB 指令
	dg.AddHandler(handler.SetStock)
	dg.AddHandler(handler.GetStock)

	// 定時排程
	handler.Task(dg)

	// 只監聽訊息
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// 開啟連線
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
