package telegrambot

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"lightft.so/WalletAutoTransfer/configuration"
)

func ConfigureChatId(config configuration.Config) {
	l := log.New(os.Stderr, "", 0)
	l.Println("log msg")

	log.Println("Initializing Bot…")
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotApiKey)
	if err != nil {
		log.Fatalf("Error initializing Bot: %s", err.Error())
	}

	tgUpdateConfig := tgbotapi.NewUpdate(0)
	tgUpdateConfig.Timeout = 30

	log.Println("Waiting for messages…")
	updates := bot.GetUpdatesChan(tgUpdateConfig)
	l.Println("\033[1mBOT READY! Please send it a message, anything will do :)\033[0m")

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msgText := fmt.Sprintf("Your chat ID is: %d", update.Message.Chat.ID)
		l.Println(msgText)
		fmt.Println(update.Message.Chat.ID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		if _, err := bot.Send(msg); err != nil {
			log.Fatalln(err.Error())
		}

		return
	}

}
