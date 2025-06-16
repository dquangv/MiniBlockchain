package main

import (
	"fmt"

	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/wallet"
)

func main() {
	// Tạo ví Alice và Bob
	alice, _ := wallet.NewWallet()
	bob, _ := wallet.NewWallet()

	fmt.Println("Alice Address:", wallet.PublicKeyToAddress(alice.PublicKey))
	fmt.Println("Bob Address:  ", wallet.PublicKeyToAddress(bob.PublicKey))

	// Tạo giao dịch
	tx := blockchain.NewTransaction(
		wallet.PublicKeyToAddress(alice.PublicKey),
		wallet.PublicKeyToAddress(bob.PublicKey),
		42.0,
	)

	tx.Sign(alice.PrivateKey)

	verified, _ := tx.Verify(alice.PublicKey)
	fmt.Println("Transaction verified:", verified)
}
