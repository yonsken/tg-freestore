package main

import (
	"log"
	"os"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := botapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := botapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
