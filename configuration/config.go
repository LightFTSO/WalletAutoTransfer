package configuration

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Network                     string `mapstructure:"NETWORK"`
	RpcUrl                      string `mapstructure:"RPC_URL"`
	ChainId                     int    `mapstructure:"CHAIN_ID"`
	Nat                         string `mapstructure:"NAT"`
	OriginWalletAddress         string `mapstructure:"ORIGIN_WALLET_ADDRESS"`
	OriginWalletPKey            string `mapstructure:"ORIGIN_WALLET_PKEY"`
	DestinationWalletAddress    string `mapstructure:"DESTINATION_WALLET_ADDRESS"`
	TelegramBotApiKey           string `mapstructure:"TELEGRAM_BOT_API_KEY"`
	TelegramBotChatId           int64  `mapstructure:"TELEGRAM_BOT_CHAT_ID"`
	TelegramNotficationsEnabled bool   `mapstructure:"TELEGRAM_NOTIFICATIONS_ENABLED"`
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

	if config.Nat == "" || config.Network == "" || config.RpcUrl == "" || config.OriginWalletAddress == "" || config.OriginWalletPKey == "" || config.DestinationWalletAddress == "" {
		log.Fatalln("Missing config values")
	}
	if config.TelegramNotficationsEnabled {
		if config.TelegramBotApiKey == "" || config.TelegramBotChatId == 0 {
			log.Fatalln("Missing Telegram config values")
		}
	}

	return
}
