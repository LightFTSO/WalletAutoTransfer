package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"lightft.so/WalletAutoTransfer/configuration"
	"lightft.so/WalletAutoTransfer/constants"
	"lightft.so/WalletAutoTransfer/functionality"
	"lightft.so/WalletAutoTransfer/telegrambot"
	"lightft.so/WalletAutoTransfer/utils"
)

func main() {
	initTelegram := flag.Bool("init-telegram", false, "Initialize Telegram Bot to get the chat ID")
	flag.Parse()

	config, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	if *initTelegram {
		telegrambot.ConfigureChatId(config)
		return
	}

	log.Println("Welcome!")
	log.Println("Connecting to", config.Network, "network. ChainId: ", config.ChainId)
	log.Println("RPC URL:", config.RpcUrl)

	var network = &constants.Network{Name: config.Network, ChainId: config.ChainId, RpcUrl: config.RpcUrl, Nat: config.Nat}

	web3Client, err := ethclient.Dial(config.RpcUrl)
	if err != nil {
		log.Fatal("Could not connect to", config.Network, "network")
		log.Fatal(err.Error())
	}

	chainId, err := web3Client.ChainID(context.Background())
	if err != nil {
		log.Fatal("Couldn't read block number")
	}
	if int(chainId.Uint64()) != config.ChainId {
		log.Fatal("The connected network's chainId (", chainId, ") doesn't match the network specified in the config.env file (", config.Network, " chainId: ", config.ChainId, ")")
	}

	/** Wallet verification section **/

	log.Println("Verifying wallets...")

	pkey, err := crypto.HexToECDSA(config.OriginWalletPKey[2:])
	if err != nil {
		log.Fatal(err.Error())
	}

	publicKey := pkey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalln("Unable to generate wallet from private key")
	}

	originAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	log.Println("Origin wallet address:", originAddress)
	if strings.ToLower(originAddress) != strings.ToLower(config.OriginWalletAddress) {
		log.Fatalln("Origin wallet address and private key don't match, please verify the wallets in the config.env file")
	}

	destAddressIsValid := utils.IsValidAddress(config.DestinationWalletAddress)
	if !destAddressIsValid {
		log.Fatalln("Destination address is not a valid address")
	}

	destinationAccount := common.HexToAddress(config.DestinationWalletAddress)

	// start telegram bot if enabled
	var tgBot telegrambot.TelegramBot
	if config.TelegramNotficationsEnabled {
		tgBot = telegrambot.TelegramBot{Bot: telegrambot.StartBot(config), ChatId: config.TelegramBotChatId, Enabled: true}
		go tgBot.SendMessage("Wallet Auto Transfer service by LightFTSO running")
		go tgBot.SendMessage(fmt.Sprintf("Connected to %s network", config.Network))
		go tgBot.SendMessage(fmt.Sprintf("Monitoring address %s", config.OriginWalletAddress))
	} else {
		tgBot = telegrambot.TelegramBot{Bot: telegrambot.StartDummyBot(), ChatId: 0, Enabled: false}
	}

	functionality.AutoTransfer(pkey, destinationAccount, web3Client, network, &tgBot)
}
