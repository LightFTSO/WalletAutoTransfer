package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"lightft.so/WalletAutoTransfer/configuration"
	"lightft.so/WalletAutoTransfer/constants"
	"lightft.so/WalletAutoTransfer/functionality"
	"lightft.so/WalletAutoTransfer/utils"
)

func Init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Println("Welcome!")

	config, err := configuration.LoadConfig("")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	network, validNetwork := constants.Networks[config.Network]
	if validNetwork == false {
		log.Fatal("The network \"", config.Network, "\" specified in the config.env file is not valid. Valid values are {Coston, Flare, Songbird}")
	}

	network.RpcUrl = config.RpcUrl

	log.Println("Connecting to", network.Name, "network")
	log.Println("RPC URL:", config.RpcUrl)

	web3Client, err := ethclient.Dial(network.RpcUrl)
	if err != nil {
		log.Fatal("Could not connect to", network.Name, "network")
		log.Fatal(err.Error())
	}

	chainId, err := web3Client.ChainID(context.Background())
	if err != nil {
		log.Fatal("Couldn't read block number")
	}
	if int(chainId.Uint64()) != network.ChainId {
		log.Fatal("The connected network's chainId (", chainId, ") doesn't match the network specified in the config.env file (", network.Name, " chainId: ", network.ChainId, ")")
	}

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

	functionality.AutoTransfer(pkey, destinationAccount, web3Client, network)
}
