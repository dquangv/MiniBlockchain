package main

import (
	"fmt"
	"log"
	"os"

	"golang-chain/pkg/p2p"
	"golang-chain/pkg/storage"
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

	nodeID := fmt.Sprintf("follower-%s", port)

	// Mở DB trước khi dùng
	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatalln("❌ Failed to open DB:", err)
	}
	defer db.Close()

	// Kiểm tra có block mới nhất không
	var localHash string
	latest, err := db.GetLatestBlock()
	if err != nil {
		fmt.Println("⚠️  No latest block found:", err)
	} else if latest != nil {
		localHash = latest.CurrentBlockHash
	}

	// Nếu là Leader → start loop tạo block định kỳ
	if peer == "none" {
		log.Println("🧠 This node is the Leader.")
		peers := []string{"localhost:50052", "localhost:50053"} // hoặc load từ config sau này
		go p2p.StartLeaderLoop(db, peers)
	} else {
		fmt.Println("🟡 Syncing block from peer:", peer)
		p2p.SyncBlocksFromPeer(peer, localHash, db)
	}

	// Khởi chạy gRPC server
	p2p.StartGRPCServer(port, dbPath, nodeID, db)
}
