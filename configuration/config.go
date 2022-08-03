package configuration

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Network                     string `mapstructure:"NETWORK"`
	RpcUrl                      string `mapstructure:"RPC_URL"`
	OriginWalletAddress         string `mapstructure:"ORIGIN_WALLET_ADDRESS"`
	OriginWalletPKey            string `mapstructure:"ORIGIN_WALLET_PKEY"`
	DestinationWalletAddress    string `mapstructure:"DESTINATION_WALLET_ADDRESS"`
	TelegramBotApiKey           string `mapstructure:"TELEGRAM_BOT_API_KEY"`
	TelegramBotChatId           int64  `mapstructure:"TELEGRAM_BOT_CHAT_ID"`
	TelegramNotficationsEnabled int    `mapstructure:"TELEGRAM_NOTIFICATIONS_ENABLED"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("$HOME/.lightftso/")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	log.Printf("Reading configuration... (%s)\n", viper.ConfigFileUsed())
	log.Println(". . .")
	if err != nil {
		log.Println(err.Error())
		panic("Unable to read in config file (check if .env file exists)")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println("Unable to parse config file (check if .env has the needed values)")
		log.Fatalln(err.Error())
	}

	return
}
