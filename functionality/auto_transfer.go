package functionality

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"lightft.so/WalletAutoTransfer/constants"
	"lightft.so/WalletAutoTransfer/telegrambot"
	"lightft.so/WalletAutoTransfer/utils"
)

/**
	In an infinite loop:

	1. Check the balance of the origin account
	2. notify the user of any errors
	3. if balance > 0, send everything - fees to destinationAccount
	4. notify the user of the transaction via the selected channel
	5. notify the user of any errors
	6. wait for the next block
**/

const blockTimesCapacity int = 500

func AutoTransfer(originPrivateKey *ecdsa.PrivateKey, destinationAccount common.Address, web3Client *ethclient.Client, network *constants.Network, tgBot *telegrambot.TelegramBot) {
	prevBlockNumber, err := web3Client.BlockNumber(context.Background())
	if err != nil {
		go tgBot.SendMessage("Error obtaining new block (1)")
		log.Fatalln(err.Error())
	}

	blockTimes := make([]uint64, 0, blockTimesCapacity)
	avgBlockTime := utils.GetAverageBlockTime(blockTimes)
	time.Sleep(time.Second)

	for {
		blockNumber, err := web3Client.BlockNumber(context.Background())

		if err != nil {
			go tgBot.SendMessage("Error obtaining new block (2)")
			log.Fatalln(err.Error())
		}
		block, err := web3Client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
		if err != nil {
			go tgBot.SendMessage("Error obtaining new block data")
			log.Fatalln(err.Error())
		}

		switch d := blockNumber - prevBlockNumber; {
		case d == 0:
			time.Sleep(time.Second)
			continue
		case d == 1:
			go CheckBalance(originPrivateKey, destinationAccount, web3Client, network.Nat, tgBot)

			prevBlockNumber = blockNumber
			blockTimes = utils.AppendToFIFOSlice(blockTimes, block.Time(), blockTimesCapacity)
			avgBlockTime = utils.GetAverageBlockTime(blockTimes)

		case d > 1:
			go CheckBalance(originPrivateKey, destinationAccount, web3Client, network.Nat, tgBot)

			missedBlocks := int(blockNumber-prevBlockNumber) - 1
			//log.Printf("Missed %d blocks\n", missedBlocks)
			for i := 1; i <= missedBlocks; i++ {
				missedBlockNumber := big.NewInt(int64(prevBlockNumber) + int64(i))
				missedBlock, err := web3Client.BlockByNumber(context.Background(), missedBlockNumber)
				if err != nil {
					go tgBot.SendMessage(fmt.Sprintf("Error obtaining new block data (2) block number: %d", missedBlockNumber))
					log.Fatalln(err.Error())
				}

				blockTimes = utils.AppendToFIFOSlice(blockTimes, missedBlock.Time(), blockTimesCapacity)
				avgBlockTime = utils.GetAverageBlockTime(blockTimes)
			}

			prevBlockNumber = blockNumber
		}

		log.Printf("Block number %d, Avg block time: %0.3fs\n", blockNumber, avgBlockTime)

		nextBlockTime := int64(block.Time()*1000 + uint64(avgBlockTime*1000))
		now := time.Now().UnixMilli()
		waitTime := time.Duration(nextBlockTime-now) * time.Millisecond
		time.Sleep(waitTime)
	}

}

func CheckBalance(originAccountPkey *ecdsa.PrivateKey, destinationAccount common.Address, web3Client *ethclient.Client, currency string, tgBot *telegrambot.TelegramBot) bool {
	_limitStr := "0.3"
	_limit, err := strconv.ParseInt(utils.ToWei(_limitStr, "ether"), 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
	}
	limit := big.NewInt(_limit)

	originAccountPublicKey := originAccountPkey.Public()
	originAccountPublicKeyECDSA, ok := originAccountPublicKey.(*ecdsa.PublicKey)
	if !ok {
		go tgBot.SendMessage("Error obtaining new block data")
		log.Fatalln("error casting public key to ECDSA")
	}
	originAccountAddress := crypto.PubkeyToAddress(*originAccountPublicKeyECDSA)

	balance, err := web3Client.BalanceAt(context.Background(), originAccountAddress, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if balance.Cmp(limit) == 1 { // if balance is greater than limit (0.3)
		msg := fmt.Sprintf("Balance is greater than %s%s, sending funds to %s", _limitStr, currency, destinationAccount)
		go tgBot.SendMessage(msg)
		log.Printf(msg)
		go sendTransaction(big.NewInt(0).Sub(balance, limit), originAccountAddress, originAccountPkey, destinationAccount, web3Client, currency, tgBot)
	}

	return true
}

func sendTransaction(value *big.Int, fromAddress common.Address, originAccountPkey *ecdsa.PrivateKey, toAddress common.Address, web3Client *ethclient.Client, currency string, tgBot *telegrambot.TelegramBot) bool {
	msg := fmt.Sprintf("Transferring %s%s from %s… to %s…\n", utils.FromWei(value.String(), "ether"), currency, fromAddress.Hex()[:10], toAddress.Hex()[:10])
	go tgBot.SendMessage(msg)
	log.Printf(msg)

	nonce, err := web3Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		go tgBot.SendMessage(fmt.Sprintf("Error obtaining nonce for address %s", fromAddress))
		log.Fatalln(err)
	}
	gasLimit := uint64(2e5)
	gasPrice, err := web3Client.SuggestGasPrice(context.Background())
	if err != nil {
		go tgBot.SendMessage("Error obtaining suggested gas price")
		log.Fatalln(err)
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := web3Client.NetworkID(context.Background())
	if err != nil {
		go tgBot.SendMessage("Error obtaining network chainId")
		log.Fatalln(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), originAccountPkey)
	if err != nil {
		go tgBot.SendMessage("Error signing transaction")
		log.Fatalln(err)
	}

	err = web3Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		go tgBot.SendMessage(fmt.Sprintf("Error sending transaction %s", signedTx.Hash().Hex()))
		log.Fatalln(err)
		return false
	}

	log.Printf("Tx hash: %s Nonce: %d", signedTx.Hash().Hex(), nonce)
	//txVerifyMSg := fmt.Sprintf("Verify it at https://coston-explorer.flare.network/tx/%s", signedTx.Hash().Hex())
	//log.Printf(txVerifyMSg)
	msg = fmt.Sprintf("Sent %s%s to %s…\n", utils.FromWei(value.String(), "ether"), currency, toAddress.Hex()[:10])
	msg2 := fmt.Sprintf("Transaction hash: %s", signedTx.Hash().Hex())
	tgBot.SendMessage(msg)
	tgBot.SendMessage(msg2)
	//tgBot.SendMessage(txVerifyMSg)
	return true
}

func notify(contents string) {

}
