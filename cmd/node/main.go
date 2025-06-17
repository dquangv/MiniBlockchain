package main

import (
	"fmt"
	"golang-chain/pkg/p2p"
	"golang-chain/pkg/storage"
	"os"
)

func main() {
	fmt.Println("Hello from validator node!")

	if len(os.Args) < 4 {
		fmt.Println("⚠️  Usage: go run main.go <port> <peer> <db-path>")
		return
	}

	port := os.Args[1]
	peer := os.Args[2]
	dbPath := os.Args[3]

	db, err := storage.NewDB(dbPath)
	if err != nil {
		fmt.Println("❌ Failed to open DB:", err)
		return
	}
	defer db.Close()

	// lấy block gần nhất nếu có
	latest, err := db.GetLatestBlock()
	if err != nil {
		fmt.Println("⚠️  No latest block found:", err)
	}
	var localHash string
	if latest != nil {
		localHash = latest.CurrentBlockHash
	}

	if peer != "none" {
		fmt.Println("🟡 Syncing block from peer:", peer)
		p2p.SyncBlocksFromPeer(peer, localHash)
	}

	p2p.StartGRPCServer(port)
}
