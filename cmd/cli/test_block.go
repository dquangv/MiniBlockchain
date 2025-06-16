package main

import (
	"fmt"
	"golang-chain/pkg/blockchain"
	"golang-chain/pkg/storage"
)

func main() {
	txs := []*blockchain.Transaction{
		blockchain.NewTransaction("alice", "bob", 5),
		blockchain.NewTransaction("bob", "alice", 2),
	}

	block := blockchain.NewBlock(txs, "genesis")

	fmt.Println("Created block with hash:", block.CurrentBlockHash)

	db, err := storage.NewDB("blockdata")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.SaveBlock(block); err != nil {
		panic(err)
	}

	loadedBlock, err := db.GetBlock(block.CurrentBlockHash)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded block Merkle Root:", loadedBlock.MerkleRoot)
}
