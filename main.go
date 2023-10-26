package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	botgolang "github.com/mail-ru-im/bot-golang"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Println("You have to create .env file")
		return
	}

	var TOKEN = os.Getenv("TOKEN")
	var API_URL = os.Getenv("API_URL")

	bot, err := botgolang.NewBot(TOKEN, botgolang.BotApiURL(API_URL), botgolang.BotDebug(true))
	if err != nil {
		log.Println("wrong token", err)
		return
	}

	ctx, finish := context.WithCancel(context.Background())
	updates := bot.GetUpdatesChannel(ctx)
	for update := range updates {
		log.Println(update)
	}

	finish()
}
