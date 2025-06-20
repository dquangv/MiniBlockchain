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

	// Nếu là Leader → start loop tạo block định kỳ
	// Xác định nodeState từ tham số truyền vào
	var nodeState p2p.NodeState
	if peer == "none" {
		nodeState = p2p.StateLeader
		log.Println("🧠 This node is the Leader.")
		peers := []string{"localhost:50052", "localhost:50053"}
		go p2p.StartLeaderLoop(db, peers)
	} else {
		nodeState = p2p.StateFollower
		fmt.Println("🟡 Syncing block from peer:", peer)
		p2p.SyncFromPeerByHeight(peer, db)
	}

	// Khởi chạy gRPC server
	p2p.StartGRPCServer(port, dbPath, nodeID, db, nodeState)
}
