package functionality

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

func AutoTransfer(originAccount common.Address, originPrivateKey *ecdsa.PrivateKey, destinationAccount common.Address, web3Client *ethclient.Client) {

	balance, err := web3Client.BalanceAt(context.Background(), originAccount, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(balance)
}

func CheckBalance(originAddress string, web3client *ethclient.Client) {

}

func mockTransaction(value *big.Int, toAddress, web3Client *ethclient.Client) {

}

func notify(contents string) {

}
