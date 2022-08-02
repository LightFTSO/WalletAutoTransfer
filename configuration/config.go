package configuration

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Network                     string `mapstructure:"NETWORK"`
	RpcUrl                      string `mapstructure:"RPC_URL"`
	DestinationWalletAddress    string `mapstructure:"DESTINATION_WALLET_ADDRESS"`
	OriginWalletAddress         string `mapstructure:"ORIGIN_WALLET_ADDRESS"`
	OriginWalletPKey            string `mapstructure:"ORIGIN_WALLET_PKEY"`
	TelegramBotApiKey           string `mapstructure:"TELEGRAM_BOT_API_KEY"`
	TelegramBotChatId           string `mapstructure:"TELEGRAM_BOT_CHAT_ID"`
	TelegramNotficationsEnabled int    `mapstructure:"TELEGRAM_NOTIFICATIONS_ENABLED"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
		panic("Unable to read in config file (check if .env file exists)")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic("Unable to parse config file (check if .env has the needed values)")
	}

	return
}
