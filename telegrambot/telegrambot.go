package telegrambot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"lightft.so/WalletAutoTransfer/configuration"
)

type TelegramBot struct {
	Bot     *tgbotapi.BotAPI
	ChatId  int64
	Enabled bool
}

func StartBot(config configuration.Config) *tgbotapi.BotAPI {
	TelegramBot, err := tgbotapi.NewBotAPI(config.TelegramBotApiKey)
	if err != nil {
		log.Fatalln(err.Error())
	}

	TelegramBot.Debug = false

	log.Printf("Authorized on account %s", TelegramBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return TelegramBot
}

func StartDummyBot() *tgbotapi.BotAPI {

	return nil
}

func (tgBot *TelegramBot) SendMessage(text string) bool {
	if !tgBot.Enabled {
		return true
	}

	msg := tgbotapi.NewMessage(tgBot.ChatId, text)
	_, err := tgBot.Bot.Send(msg)
	if err != nil {
		log.Fatalln("Unable to send telegram notification")
	}
	return true
}
