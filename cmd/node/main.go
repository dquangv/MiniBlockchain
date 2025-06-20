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

	if len(os.Args) < 3 {
		fmt.Println("⚠️  Usage: go run main.go <port> <db-path>")
		return
	}

	port := os.Args[1]
	dbPath := os.Args[2]
	allPeers := []string{"localhost:50051", "localhost:50052", "localhost:50053"}

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
	var peers []string

	leader := p2p.DetectLeader(allPeers)
	if leader == "" || port == "50051" { // nếu không thấy leader hoặc là node 50051 thì làm Leader
		nodeState = p2p.StateLeader
		log.Println("🧠 This node is the Leader.")
		peers = []string{"localhost:50052", "localhost:50053"}
		go p2p.StartLeaderLoop(db, peers)
	} else {
		nodeState = p2p.StateFollower
		log.Println("🔄 Detected Leader at", leader, "→ syncing...")
		p2p.SyncFromPeerByHeight(leader, db)
	}

	// Khởi chạy gRPC server
	p2p.StartGRPCServer(port, dbPath, nodeID, db, nodeState)
}
