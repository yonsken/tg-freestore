package main

import (
	"log"
	"os"
	"path/filepath"

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
		if err := processUpdate(bot, update); err != nil {
			panic(err)
		}
	}
}

func processUpdate(bot *botapi.BotAPI, update botapi.Update) error {
	if update.Message == nil {
		return nil
	}

	var (
		userName = update.Message.From.UserName
		chatID   = update.FromChat().ChatConfig().ChatID
	)

	log.Printf("Update message from [%s]: %s", userName, update.Message.Text)

	fileReader, err := os.Open(filepath.Clean("test/gopher.png"))
	if err != nil {
		return err
	}

	file := botapi.FileReader{
		Name:   "gopher.png",
		Reader: fileReader,
	}

	mediaGroup := botapi.NewMediaGroup(chatID, []any{botapi.NewInputMediaPhoto(file)})

	messages, err := bot.SendMediaGroup(mediaGroup)
	if err != nil {
		return err
	}

	log.Printf("Successfully sent test file to chat %d with [%s]", chatID, userName)

	for i, msg := range messages {
		log.Printf("Send Media Group message %d: %s", i, msg.Text)
	}

	return nil
}
